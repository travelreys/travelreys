export interface Op {
  op: "add" | "remove" | "replace"
  path: string
  value: string
};

/**
 * Should avoid replacing entire list
 */
export const makeRepOp = (path: string, value: any): Op => {
  return { op: "replace", path, value }
};

export const makeAddOp = (path: string, value: any): Op => {
  return { op: "add", path, value }
};

export const makeRemoveOp = (path: string, value: any): Op => {
  return { op: "remove", path, value }
};

export const jsonPathToPath = (jsonPath: string): string => {
  return jsonPath.replace("/", ".")
}
