export interface TemplateVariables {
  id: string;
  template_id: string;
  key: string;
  default_value: string;
  required: boolean;
  description: string;
}

export type TemplateEditVariable = {
  key: string;
  default_value: string;
  description: string;
  required : boolean
};