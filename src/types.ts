import { DataQuery, DataSourceJsonData } from '@grafana/data';
import { FORMAT_TIME_SERIES } from 'DataSource';
export interface PrestoQuery extends DataQuery {
  rawSql: string;
  format: string;
  legendFormat: string;
  queryText?: string;
  queryType?: string;
}

export const defaultQuery: Partial<PrestoQuery> = {
  rawSql: '',
  format: FORMAT_TIME_SERIES,
};

/**
 * These are options configured for each DataSource instance
 */
export interface PrestoDataSourceOptions extends DataSourceJsonData {
  host?: string;
  catalog?: string;
  schema?: string;
  queryMaxExecutionSeconds?: number;
  rowLimit?: number;
  resultRowLimit?: number;
  customParams: CustomParam[];
}

export interface CustomParam {
  name: string;
  value: string;
}
