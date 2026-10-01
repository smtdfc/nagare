import { Link } from "@tanstack/react-router";
import { ChevronRight } from "lucide-react";

export interface SettingItem {
  title: string;
  description: string;
  icon: React.ElementType;
  href: string;
}

export interface SettingSection {
  category: string;
  items: SettingItem[];
}

export function SettingItem({ item }: { item: SettingItem }) {
  const Icon = item.icon;
  return (
    <Link
      to={item.href as any}
      className="flex items-center justify-between p-4 rounded-xl border bg-card hover:bg-accent/50 transition-all duration-200 group shadow-sm"
    >
      <div className="flex items-center gap-4">
        <div className="p-2.5 rounded-lg bg-primary/10 text-primary group-hover:bg-primary group-hover:text-primary-foreground transition-colors">
          <Icon className="w-5 h-5" />
        </div>
        <div>
          <h5 className="text-md font-medium text-card-foreground">
            {item.title}
          </h5>
          <p className="text-xs text-muted-foreground">{item.description}</p>
        </div>
      </div>
      <ChevronRight className="w-5 h-5 text-muted-foreground group-hover:translate-x-0.5 transition-transform" />
    </Link>
  );
}
