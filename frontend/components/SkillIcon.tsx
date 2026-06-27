import {
  Hand,
  Sparkles,
  Coffee,
  ShoppingBag,
  Plane,
  Users,
  MessageCircle,
  Heart,
  Compass,
  Hash,
  BookOpen,
  type LucideIcon,
} from "lucide-react";

// Maps the lucide icon names stored on a Skill to their components, so skill
// artwork stays clean and professional (no emoji).
const ICONS: Record<string, LucideIcon> = {
  Hand,
  Sparkles,
  Coffee,
  ShoppingBag,
  Plane,
  Users,
  MessageCircle,
  Heart,
  Compass,
  Hash,
  BookOpen,
};

export function SkillIcon({
  name,
  size = 24,
  className = "",
}: {
  name?: string;
  size?: number;
  className?: string;
}) {
  const Icon = (name && ICONS[name]) || BookOpen;
  return <Icon size={size} className={className} strokeWidth={2} />;
}
