const state = {
  token: localStorage.getItem("hmk_access") || "",
  refresh: localStorage.getItem("hmk_refresh") || "",
};

const authScreen = document.getElementById("auth-screen");
const app = document.getElementById("app");
const view = document.getElementById("view");
const title = document.getElementById("view-title");
const authError = document.getElementById("auth-error");

function headers() {
  return {
    "Content-Type": "application/json",
    Authorization: state.token ? `Bearer ${state.token}` : "",
  };
}

async function api(path, options = {}) {
  const res = await fetch(path, {
    ...options,
    headers: { ...headers(), ...(options.headers || {}) },
  });
  const body = await res.json().catch(() => ({}));
  if (!res.ok) throw new Error(body.error || body.message || res.statusText);
  return body;
}

function showApp() {
  authScreen.classList.add("hidden");
  app.classList.remove("hidden");
  render("users");
}

function showAuth() {
  app.classList.add("hidden");
  authScreen.classList.remove("hidden");
}

document.getElementById("login-form").addEventListener("submit", async (e) => {
  e.preventDefault();
  authError.textContent = "";
  const data = Object.fromEntries(new FormData(e.target));
  try {
    const body = await api("/api/v1/auth/login", { method: "POST", body: JSON.stringify(data) });
    state.token = body.data.access_token;
    state.refresh = body.data.refresh_token;
    localStorage.setItem("hmk_access", state.token);
    localStorage.setItem("hmk_refresh", state.refresh);
    showApp();
  } catch (err) {
    authError.textContent = err.message;
  }
});

document.getElementById("bootstrap-form").addEventListener("submit", async (e) => {
  e.preventDefault();
  authError.textContent = "";
  const data = Object.fromEntries(new FormData(e.target));
  try {
    await api("/api/v1/auth/bootstrap", { method: "POST", body: JSON.stringify(data) });
    authError.textContent = "Plateforme créée. Connectez-vous.";
  } catch (err) {
    authError.textContent = err.message;
  }
});

document.getElementById("logout").addEventListener("click", () => {
  localStorage.removeItem("hmk_access");
  localStorage.removeItem("hmk_refresh");
  state.token = "";
  showAuth();
});

document.querySelectorAll("aside nav button").forEach((btn) => {
  btn.addEventListener("click", () => {
    document.querySelectorAll("aside nav button").forEach((b) => b.classList.remove("active"));
    btn.classList.add("active");
    render(btn.dataset.view);
  });
});

function table(rows, columns) {
  if (!rows || !rows.length) return "<p class='muted'>Aucune donnée.</p>";
  return `<table><thead><tr>${columns.map((c) => `<th>${c}</th>`).join("")}</tr></thead><tbody>
    ${rows.map((row) => `<tr>${columns.map((c) => `<td>${escapeHtml(fmt(row[c]))}</td>`).join("")}</tr>`).join("")}
  </tbody></table>`;
}

function fmt(value) {
  if (value == null) return "";
  if (Array.isArray(value)) return value.join(", ");
  if (typeof value === "object") return JSON.stringify(value);
  return String(value);
}

function escapeHtml(value) {
  return value.replace(/[&<>"']/g, (ch) => ({ "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#39;" }[ch]));
}

async function render(name) {
  const titles = {
    users: "Utilisateurs",
    tenants: "Tenants",
    roles: "Rôles",
    permissions: "Permissions",
    keys: "API Keys",
    secrets: "Secrets",
    audit: "Audit",
    fraud: "Règles anti-fraude",
    settings: "Paramètres globaux",
  };
  title.textContent = titles[name] || name;
  try {
    if (name === "users") {
      const body = await api("/api/v1/users");
      view.innerHTML = panelForm("Créer un utilisateur", [
        ["email", "email"], ["name", "text"], ["password", "password"], ["status", "text"],
      ], "/api/v1/users") + table(body.data.users, ["id", "email", "name", "status"]);
    } else if (name === "tenants") {
      const body = await api("/api/v1/tenants");
      view.innerHTML = panelForm("Créer un tenant", [
        ["name", "text"], ["slug", "text"], ["plan", "text"], ["status", "text"],
      ], "/api/v1/tenants") + table(body.data.tenants, ["id", "name", "slug", "plan", "status"]);
    } else if (name === "roles") {
      const body = await api("/api/v1/rbac/roles");
      view.innerHTML = panelForm("Créer un rôle", [
        ["name", "text"], ["description", "text"],
      ], "/api/v1/rbac/roles") + table(body.data.roles, ["id", "name", "description", "permissions"]);
    } else if (name === "permissions") {
      const body = await api("/api/v1/permissions");
      view.innerHTML = panelForm("Créer une permission", [
        ["name", "text"], ["resource", "text"], ["action", "text"], ["description", "text"],
      ], "/api/v1/permissions") + table(body.data.permissions, ["id", "name", "resource", "action"]);
    } else if (name === "keys") {
      const body = await api("/api/v1/api-keys");
      view.innerHTML = panelForm("Créer une API key", [["name", "text"]], "/api/v1/api-keys") + table(body.data.api_keys, ["id", "name", "prefix", "revoked_at", "expires_at"]);
    } else if (name === "secrets") {
      const body = await api("/api/v1/secrets");
      view.innerHTML = panelForm("Créer un secret", [
        ["name", "text"], ["type", "text"], ["value", "password"], ["description", "text"],
      ], "/api/v1/secrets") + table(body.data.secrets, ["id", "name", "type", "masked_value", "revoked_at"]);
    } else if (name === "audit") {
      const body = await api("/api/v1/audit/events");
      view.innerHTML = table(body.data.events, ["created_at", "action", "resource", "user_id", "status"]);
    } else if (name === "fraud") {
      const body = await api("/api/v1/fraud/rules");
      view.innerHTML = panelForm("Créer une règle", [
        ["name", "text"], ["field", "text"], ["operator", "text"], ["threshold", "text"], ["action", "text"], ["score", "number"],
      ], "/api/v1/fraud/rules") + table(body.data.rules, ["name", "field", "operator", "threshold", "action", "enabled"]);
    } else if (name === "settings") {
      const body = await api("/api/v1/settings");
      view.innerHTML = panelForm("Enregistrer un paramètre", [
        ["key", "text"], ["value", "text"],
      ], "/api/v1/settings", "POST", (data) => {
        let value = data.value;
        try { value = JSON.parse(value); } catch (_) {}
        return { key: data.key, value };
      }) + table(body.data.settings, ["key", "value", "updated_at"]);
    }
    bindForms();
  } catch (err) {
    view.innerHTML = `<p class="error">${escapeHtml(err.message)}</p>`;
  }
}

function panelForm(titleText, fields, path, method = "POST", transform) {
  const inputs = fields.map(([name, type]) => `<label>${name}<input name="${name}" type="${type}" /></label>`).join("");
  return `<div class="panel"><h3>${titleText}</h3><form data-path="${path}" data-method="${method}">${inputs}<button type="submit">Enregistrer</button></form></div>`;
}

function bindForms() {
  view.querySelectorAll("form").forEach((form) => {
    form.addEventListener("submit", async (e) => {
      e.preventDefault();
      const data = Object.fromEntries(new FormData(form));
      if (data.score) data.score = Number(data.score);
      if (data.enabled === undefined && form.dataset.path.includes("fraud")) data.enabled = true;
      try {
        await api(form.dataset.path, { method: form.dataset.method || "POST", body: JSON.stringify(data) });
        const current = document.querySelector("aside nav button.active")?.dataset.view || "users";
        render(current);
      } catch (err) {
        alert(err.message);
      }
    });
  });
}

if (state.token) showApp();
