export type Templates = {
  category: string;
  created_at: string;
  description: string;
  id: string;
  image: string;
  is_public: boolean;
  name: string;
  template_url: string;
  icon: string;
};

export type TemplateForm = {
  template_name: string;
  variables: {
    key: string;
    required: boolean;
  }[];
  addons: {
    id: string;
    name: string;
  }[];
};

export type TemplateVariable = {
  key: string;
  default_value: string;
  required: boolean;
  description: string;
};

export type TemplateAddon = {
  name: string;
  image: string;
  description: string;
  default_config: Record<string, unknown>;
};

export type TemplateFiles = {
  filename: string;
  sort_order: number;
};

export type CreateTemplateRequest = {
  name: string;
  description: string;
  image: string;
  category: string;
  icon: string;
  is_public: boolean;
  variables: TemplateVariable[];
  addons: TemplateAddon[];
  files: TemplateFiles[];
};

export interface TemplateAddOn {
  id: string;
  template_id: string;
  name: string;
  image: string;
  description: string;
}

export type EditState = {
  name: string;
  image: string;
  description: string;
};
