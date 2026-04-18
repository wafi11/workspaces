export interface AwarenessUser {
  clientId: number;
  name: string;
  color: string;
  colorLight: string;
}

export interface CollabEditorProps {
  workspaceId: string;
  filePath?: string;
}