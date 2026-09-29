import { apiDelete, apiGet, apiPost, apiPut, normalizeList, type ApiResponse } from './api';

export interface FraudRuleRecord {
  id?: string;
  name?: string;
  enabled?: boolean;
  field?: string;
  operator?: string;
  threshold?: string;
  action?: string;
  score?: number;
  tenant_id?: string;
  window_seconds?: number;
  metadata?: Record<string, unknown>;
  created_at?: string;
  updated_at?: string;
}

export interface FraudViolationRecord {
  id?: string;
  tenant_id?: string;
  user_id?: string;
  signal_id?: string;
  rule_id?: string;
  score?: number;
  level?: string;
  reason?: string;
  status?: string;
  metadata?: Record<string, unknown>;
  created_at?: string;
  updated_at?: string;
}

export interface FraudRuleWrite {
  id?: string;
  name: string;
  enabled: boolean;
  field: string;
  operator: string;
  threshold: string;
  window_seconds?: number;
  action: string;
  score: number;
  metadata?: Record<string, unknown>;
}

export interface FraudAnalysisResult {
  risk_score: number;
  risk_level: string;
  blocked: boolean;
  matched: string[];
  analyzed_at: string;
}

export async function fetchFraudRules(): Promise<FraudRuleRecord[]> {
  const payload = await apiGet<unknown>('/api/v1/fraud/rules');
  return normalizeList<FraudRuleRecord>(payload);
}

export async function createFraudRule(rule: FraudRuleWrite): Promise<void> {
  await apiPost('/api/v1/fraud/rules', rule);
}

export async function updateFraudRule(rule: FraudRuleWrite & { id: string }): Promise<void> {
  await apiPut('/api/v1/fraud/rules', rule);
}

export async function deleteFraudRule(id: string): Promise<void> {
  await apiDelete('/api/v1/fraud/rules', { params: { id } });
}

export async function analyzeFraud(attributes: Record<string, unknown>): Promise<ApiResponse<FraudAnalysisResult>> {
  return apiPost<ApiResponse<FraudAnalysisResult>>('/api/v1/fraud/analyze', attributes);
}

export async function fetchFraudViolations(): Promise<FraudViolationRecord[]> {
  const payload = await apiGet<unknown>('/api/v1/fraud/violations');
  return normalizeList<FraudViolationRecord>(payload);
}
