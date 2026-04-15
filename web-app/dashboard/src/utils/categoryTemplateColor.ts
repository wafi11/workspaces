export const getCategoryStyle = (category: string) => {
  const cat = category.toLowerCase();
  if (cat.includes('editor')) return 'bg-blue-500/10 text-blue-400 ring-blue-500/20';
  if (cat.includes('database')) return 'bg-emerald-500/10 text-emerald-400 ring-emerald-500/20';
//   if (cat.includes('')) return 'bg-purple-500/10 text-purple-400 ring-purple-500/20';
  return 'bg-sidebar-border/50 text-sidebar-text-active ring-sidebar-border';
};