import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Alert } from '../../components/ui/Alert';
import { Badge } from '../../components/ui/Badge';
import { Button } from '../../components/ui/Button';
import { Card } from '../../components/ui/Card';

const capabilityGroups = [
  {
    eyebrow: 'MULTI-TENANT FOUNDATION',
    title: 'One security control plane for every organization you operate.',
    description:
      'HostMK Security Framework gives each organization a distinct operational boundary while keeping governance, visibility, and policy management in one enterprise-grade platform.',
    points: [
      'Tenant-scoped identities, sessions, API keys, secrets, rules, and audit records.',
      'Organization lifecycle controls for activation, suspension, ownership, and operational continuity.',
      'A scalable model for SaaS platforms, managed service providers, and multi-client enterprise teams.',
    ],
  },
  {
    eyebrow: 'IDENTITY & GOVERNANCE',
    title: 'Make access deliberate, reviewable, and accountable.',
    description:
      'Identity capabilities are designed for teams that need to manage people, responsibilities, and privileged operations without losing tenant boundaries or an evidence trail.',
    points: [
      'User lifecycle management with activation, suspension, verification, MFA enrollment, and role assignments.',
      'Enterprise RBAC with granular permissions and a clear separation between roles and effective access.',
      'Session visibility and revocation for rapid response to compromised devices or changed responsibilities.',
    ],
  },
  {
    eyebrow: 'PLATFORM SECURITY',
    title: 'Operate integrations and sensitive material with stronger controls.',
    description:
      'API management and secret operations are built around tenant ownership, controlled access, lifecycle management, and auditable change paths.',
    points: [
      'API key inventory, scoped access, lifecycle awareness, rotation, and revocation.',
      'Encrypted secret storage with tenant context, controlled rotation, revocation, and version history.',
      'Authentication, token refresh, rate controls, security headers, and defense-in-depth middleware in the backend platform.',
    ],
  },
];

const operationalCapabilities = [
  {
    title: 'Audit & compliance',
    body: 'Create a durable trail around operational changes, security events, actor context, resource scope, and outcomes. Investigations have context instead of fragmented logs.',
  },
  {
    title: 'Fraud detection & risk analysis',
    body: 'Define tenant-scoped rules, evaluate attributes, record violations, score risk, and identify actions that require additional scrutiny or blocking.',
  },
  {
    title: 'Monitoring & observability',
    body: 'Expose application health, database reachability, service identity, version, and environment through a focused operational monitoring surface.',
  },
  {
    title: 'Performance & scalability',
    body: 'A focused Go service, tenant-aware persistence, rate controls, and explicit lifecycle operations provide a clear path for dependable growth.',
  },
];

const useCases = [
  {
    title: 'B2B SaaS platforms',
    body: 'Give every customer organization isolated users, roles, sessions, audit records, API access, and security controls without operating separate security stacks.',
  },
  {
    title: 'Managed service providers',
    body: 'Operate multiple client environments with consistent governance while keeping credentials, policies, evidence, and operational boundaries distinct.',
  },
  {
    title: 'Enterprise internal platforms',
    body: 'Centralize access governance for internal tools, platform teams, business units, and integration workloads that need clear accountability.',
  },
  {
    title: 'Multi-client digital products',
    body: 'Build a trustworthy security foundation for products that must scale organizations, integrate external systems, and respond to evolving risk.',
  },
];

const roadmap = [
  ['Available foundation', 'Tenant operations, user management, RBAC, permissions, session controls, API keys, secrets, audit trails, fraud rules, settings, and health monitoring are present in the platform.'],
  ['Next integration phase', 'A dedicated authentication interface will connect enterprise sign-in journeys to the existing backend authentication contract.'],
  ['Operational maturity', 'The product roadmap continues toward richer dashboards, workflow-driven administration, and deeper policy visibility without compromising tenant isolation.'],
];

const faq = [
  ['How does multi-tenant isolation work?', 'Platform resources are scoped to a tenant context. Users, sessions, roles, permissions, secrets, API keys, audit events, and fraud data are handled through tenant-aware backend operations.'],
  ['How are permissions managed?', 'Roles are associated with granular permissions. The platform evaluates access through its authorization middleware before protected API resources are served.'],
  ['How are secrets protected?', 'Secrets are stored through the backend secret-management capability, which uses encryption, tenant ownership, rotation, revocation, and version-history operations.'],
  ['What does the audit trail capture?', 'Audit records capture security-relevant actions with actor, tenant, resource, message, status, severity, network, and user-agent context where the operation provides it.'],
  ['How does fraud detection work?', 'Tenant-scoped fraud rules evaluate supplied attributes. Matching rules create violations, contribute to risk scoring, and can trigger a blocking outcome.'],
  ['How can an application integrate with the platform?', 'The backend exposes versioned HTTP endpoints for managed resources. API keys, permissions, rate controls, and audit capabilities support governed integrations.'],
  ['Can one team manage several organizations?', 'Yes. The tenant management capability is designed to create, list, activate, suspend, and govern organization boundaries from a central control plane.'],
  ['How are API keys rotated?', 'The API key resource supports a dedicated rotation action and lifecycle controls including revocation. Plaintext handling remains an explicit, controlled backend concern.'],
  ['How are sessions managed?', 'The platform records sessions in the tenant context and provides an explicit revocation action for ending access when a device or account must be contained.'],
  ['How is platform health supervised?', 'The health endpoint reports service status and component health, including database reachability. The version endpoint provides application identity and environment context.'],
  ['Does the platform support MFA?', 'The user-management backend includes MFA enrollment and validates MFA during authentication for users whose MFA is enabled.'],
  ['What is the role of the admin console?', 'The console gives operators a focused interface for tenant, identity, session, RBAC, API, secret, audit, fraud, settings, and monitoring workflows.'],
];

export function LandingPage() {
  const navigate = useNavigate();
  const [loginNoticeVisible, setLoginNoticeVisible] = useState(false);

  return (
    <div className="min-h-screen bg-slate-950 text-slate-100">
      <a href="#main-content" className="sr-only focus:not-sr-only focus:absolute focus:left-4 focus:top-4 focus:z-50 focus:rounded-md focus:bg-sky-400 focus:px-4 focus:py-2 focus:text-slate-950">
        Skip to content
      </a>

      <header className="sticky top-0 z-40 border-b border-slate-800 bg-slate-950/95 backdrop-blur">
        <nav aria-label="Primary navigation" className="mx-auto flex max-w-7xl items-center justify-between gap-4 px-4 py-4 sm:px-6 lg:px-8">
          <a href="#top" className="min-w-0 text-left">
            <span className="block text-xs font-semibold uppercase tracking-[0.22em] text-sky-400">HostMK</span>
            <span className="block text-base font-semibold text-white">Security Framework</span>
          </a>
          <div className="hidden items-center gap-5 text-sm text-slate-300 lg:flex">
            <a className="hover:text-sky-300" href="#platform">Platform</a>
            <a className="hover:text-sky-300" href="#security">Security</a>
            <a className="hover:text-sky-300" href="#operations">Operations</a>
            <a className="hover:text-sky-300" href="#faq">FAQ</a>
          </div>
          <Button type="button" onClick={() => navigate('/')}>Open dashboard</Button>
        </nav>
      </header>

      <main id="main-content">
        <section id="top" className="relative overflow-hidden border-b border-slate-800">
          <div className="mx-auto grid max-w-7xl gap-12 px-4 py-20 sm:px-6 lg:grid-cols-[1.15fr_0.85fr] lg:px-8 lg:py-28">
            <div className="max-w-3xl">
              <Badge className="border-sky-500/40 bg-sky-500/10 text-sky-200">Enterprise security operations for multi-tenant platforms</Badge>
              <h1 className="mt-6 text-4xl font-semibold tracking-tight text-white sm:text-5xl lg:text-6xl">Security governance that scales with every organization you serve.</h1>
              <p className="mt-6 max-w-2xl text-lg leading-8 text-slate-300">HostMK Security Framework unifies tenant isolation, identity, RBAC, API access, secret operations, audit evidence, fraud controls, and service health into one disciplined platform foundation.</p>
              <div className="mt-8 flex flex-col gap-3 sm:flex-row">
                <Button type="button" className="border-sky-400 bg-sky-400 text-slate-950 hover:bg-sky-300" onClick={() => navigate('/')}>Explore the dashboard</Button>
                <Button type="button" className="bg-transparent" onClick={() => setLoginNoticeVisible((visible) => !visible)} aria-expanded={loginNoticeVisible}>Login availability</Button>
              </div>
              {loginNoticeVisible ? <Alert className="mt-4 border-amber-500/40 bg-amber-950/40 text-amber-100" role="status">The backend authentication contract is available. A dedicated enterprise sign-in interface will be introduced in a separate phase; no temporary login route is exposed.</Alert> : null}
              <div className="mt-10 grid grid-cols-1 gap-3 sm:grid-cols-3">
                {['Tenant-aware boundaries', 'Policy-driven access', 'Evidence-ready operations'].map((item) => <div key={item} className="border-l-2 border-sky-400 pl-3 text-sm font-medium text-slate-200">{item}</div>)}
              </div>
            </div>
            <Card className="border-slate-700 bg-slate-900/90 p-6 shadow-2xl shadow-sky-950/40">
              <p className="text-sm font-semibold uppercase tracking-[0.18em] text-sky-300">Operational assurance</p>
              <div className="mt-6 space-y-5">
                <div><p className="text-2xl font-semibold text-white">Tenant context</p><p className="mt-1 text-sm text-slate-400">Security resources and operational decisions remain tied to the organization they serve.</p></div>
                <div className="border-t border-slate-800 pt-5"><p className="text-2xl font-semibold text-white">Governed access</p><p className="mt-1 text-sm text-slate-400">Roles, permissions, API keys, sessions, and authentication controls work together instead of in isolation.</p></div>
                <div className="border-t border-slate-800 pt-5"><p className="text-2xl font-semibold text-white">Traceable outcomes</p><p className="mt-1 text-sm text-slate-400">Auditing, fraud events, health checks, and lifecycle actions support informed operational decisions.</p></div>
              </div>
            </Card>
          </div>
        </section>

        <section id="platform" className="mx-auto max-w-7xl px-4 py-20 sm:px-6 lg:px-8">
          <div className="max-w-3xl"><p className="text-sm font-semibold uppercase tracking-[0.18em] text-sky-400">Platform foundation</p><h2 className="mt-3 text-3xl font-semibold text-white sm:text-4xl">A practical security operating layer for products with real organizational complexity.</h2><p className="mt-5 text-lg leading-8 text-slate-300">HostMK Security Framework is built for teams that must protect multiple organizations, govern privileged operations, and keep evidence close to the work being performed. It reduces fragmented identity and security workflows without flattening important tenant boundaries.</p></div>
          <div className="mt-10 grid gap-6 lg:grid-cols-3">
            {capabilityGroups.map((group) => <Card key={group.eyebrow} className="h-full"><p className="text-xs font-semibold tracking-[0.16em] text-sky-400">{group.eyebrow}</p><h3 className="mt-3 text-xl font-semibold text-white">{group.title}</h3><p className="mt-3 text-sm leading-6 text-slate-300">{group.description}</p><ul className="mt-5 space-y-3 text-sm leading-6 text-slate-400">{group.points.map((point) => <li key={point} className="flex gap-3"><span aria-hidden="true" className="text-sky-400">•</span><span>{point}</span></li>)}</ul></Card>)}
          </div>
        </section>

        <section id="security" className="border-y border-slate-800 bg-slate-900/40"><div className="mx-auto max-w-7xl px-4 py-20 sm:px-6 lg:px-8"><div className="grid gap-10 lg:grid-cols-[0.9fr_1.1fr]"><div><p className="text-sm font-semibold uppercase tracking-[0.18em] text-sky-400">Enterprise controls</p><h2 className="mt-3 text-3xl font-semibold text-white">Build trustworthy access paths from user identity to sensitive operations.</h2><p className="mt-5 leading-7 text-slate-300">The framework combines authentication, MFA support, encrypted secret handling, authorization middleware, API key controls, rate limiting, security headers, and tenant-aware persistence. Each layer contributes to a defensible platform posture.</p></div><div className="grid gap-4 sm:grid-cols-2">{['User and organization management', 'RBAC enterprise governance', 'API management and key lifecycle', 'Encrypted secrets management', 'Audit and compliance evidence', 'Fraud detection and risk scoring', 'Session control and revocation', 'Defense-in-depth backend controls'].map((item) => <div key={item} className="rounded-lg border border-slate-700 bg-slate-950/70 p-4 text-sm font-medium text-slate-200">{item}</div>)}</div></div></div></section>

        <section id="operations" className="mx-auto max-w-7xl px-4 py-20 sm:px-6 lg:px-8"><div className="text-center"><p className="text-sm font-semibold uppercase tracking-[0.18em] text-sky-400">Operational intelligence</p><h2 className="mt-3 text-3xl font-semibold text-white">Security is most useful when it improves the quality of operations.</h2></div><div className="mt-10 grid gap-6 md:grid-cols-2">{operationalCapabilities.map((capability) => <Card key={capability.title}><h3 className="text-xl font-semibold text-white">{capability.title}</h3><p className="mt-3 leading-7 text-slate-300">{capability.body}</p></Card>)}</div></section>

        <section className="border-y border-slate-800 bg-slate-900/40"><div className="mx-auto max-w-7xl px-4 py-20 sm:px-6 lg:px-8"><div className="grid gap-10 lg:grid-cols-2"><div><p className="text-sm font-semibold uppercase tracking-[0.18em] text-sky-400">Built for real teams</p><h2 className="mt-3 text-3xl font-semibold text-white">One platform foundation, many secure operating models.</h2><p className="mt-5 leading-7 text-slate-300">Use HostMK Security Framework where growth, accountability, integration, and customer trust must advance together.</p></div><div className="grid gap-4">{useCases.map((useCase) => <div key={useCase.title} className="border-l-2 border-sky-400 bg-slate-950/40 p-4"><h3 className="font-semibold text-white">{useCase.title}</h3><p className="mt-2 text-sm leading-6 text-slate-300">{useCase.body}</p></div>)}</div></div></div></section>

        <section className="mx-auto max-w-7xl px-4 py-20 sm:px-6 lg:px-8"><div className="grid gap-10 lg:grid-cols-[0.8fr_1.2fr]"><div><p className="text-sm font-semibold uppercase tracking-[0.18em] text-sky-400">Product roadmap</p><h2 className="mt-3 text-3xl font-semibold text-white">Progressive capability, grounded in the platform that already exists.</h2><p className="mt-5 leading-7 text-slate-300">The roadmap prioritizes the operational surfaces that turn strong backend controls into a coherent enterprise experience.</p></div><ol className="space-y-5">{roadmap.map(([title, body], index) => <li key={title} className="flex gap-4"><span className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full border border-sky-500/50 bg-sky-500/10 text-sm font-semibold text-sky-200">{index + 1}</span><div><h3 className="font-semibold text-white">{title}</h3><p className="mt-1 leading-7 text-slate-300">{body}</p></div></li>)}</ol></div></section>

        <section id="faq" className="border-t border-slate-800 bg-slate-900/40"><div className="mx-auto max-w-5xl px-4 py-20 sm:px-6 lg:px-8"><p className="text-center text-sm font-semibold uppercase tracking-[0.18em] text-sky-400">Professional FAQ</p><h2 className="mt-3 text-center text-3xl font-semibold text-white">Answers for teams evaluating a security foundation.</h2><div className="mt-10 divide-y divide-slate-800 rounded-xl border border-slate-800 bg-slate-950/60">{faq.map(([question, answer]) => <details key={question} className="group p-5"><summary className="cursor-pointer list-none pr-8 font-medium text-white marker:hidden focus:outline-none focus:ring-2 focus:ring-sky-500"><span>{question}</span><span aria-hidden="true" className="float-right text-sky-400 group-open:rotate-45">+</span></summary><p className="mt-4 max-w-3xl leading-7 text-slate-300">{answer}</p></details>)}</div></div></section>

        <section className="border-t border-slate-800"><div className="mx-auto max-w-7xl px-4 py-16 text-center sm:px-6 lg:px-8"><Badge className="border-sky-500/40 bg-sky-500/10 text-sky-200">A disciplined foundation for secure growth</Badge><h2 className="mx-auto mt-5 max-w-3xl text-3xl font-semibold text-white">Bring tenant operations, security governance, and product resilience into the same conversation.</h2><p className="mx-auto mt-4 max-w-2xl leading-7 text-slate-300">Open the existing admin dashboard to explore the operational capabilities already available in the platform.</p><div className="mt-7 flex justify-center"><Button type="button" className="border-sky-400 bg-sky-400 text-slate-950 hover:bg-sky-300" onClick={() => navigate('/')}>Go to dashboard</Button></div></div></section>
      </main>

      <footer className="border-t border-slate-800 bg-slate-950"><div className="mx-auto grid max-w-7xl gap-10 px-4 py-12 sm:px-6 md:grid-cols-2 lg:grid-cols-4 lg:px-8"><div><p className="text-sm font-semibold uppercase tracking-[0.18em] text-sky-400">HostMK</p><p className="mt-2 font-semibold text-white">Security Framework</p><p className="mt-3 text-sm leading-6 text-slate-400">Enterprise security operations for the organizations your platform serves.</p></div><div><h2 className="text-sm font-semibold text-white">Platform</h2><ul className="mt-3 space-y-2 text-sm text-slate-400"><li><a href="#platform" className="hover:text-sky-300">Multi-tenant architecture</a></li><li><a href="#security" className="hover:text-sky-300">Access governance</a></li><li><a href="#operations" className="hover:text-sky-300">Operational intelligence</a></li></ul></div><div><h2 className="text-sm font-semibold text-white">Resources</h2><ul className="mt-3 space-y-2 text-sm text-slate-400"><li><a href="#faq" className="hover:text-sky-300">FAQ</a></li><li><button type="button" className="text-left hover:text-sky-300" onClick={() => setLoginNoticeVisible(true)}>Authentication interface</button></li><li><button type="button" className="text-left hover:text-sky-300" onClick={() => navigate('/')}>Admin dashboard</button></li></ul></div><div><h2 className="text-sm font-semibold text-white">Trust</h2><ul className="mt-3 space-y-2 text-sm text-slate-400"><li>Tenant isolation</li><li>Auditability and traceability</li><li>Secrets and key lifecycle controls</li><li>Role-based governance</li></ul></div></div><div className="border-t border-slate-800 px-4 py-6 text-center text-sm text-slate-500 sm:px-6 lg:px-8">© {new Date().getFullYear()} HostMK Security Framework. Security architecture for modern multi-tenant platforms.</div></footer>
    </div>
  );
}
