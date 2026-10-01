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
  GraduationCap,
  PenLine,
  Quote,
  Languages,
  Layers,
  Clock,
  Link2,
  Music,
  Waves,
  CheckCircle,
  Globe,
  Scale,
  Link,
  Shield,
  HeartPulse,
  Briefcase,
  Home,
  Smartphone,
  Landmark,
  FlaskConical,
  Brain,
  Cpu,
  Leaf,
  Mic,
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
  GraduationCap,
  PenLine,
  Quote,
  Languages,
  Layers,
  Clock,
  Link2,
  // Mandarin course
  Music,
  Waves,
  CheckCircle,
  Globe,
  Scale,
  Link,
  Shield,
  HeartPulse,
  Briefcase,
  Home,
  Smartphone,
  Landmark,
  FlaskConical,
  Brain,
  Cpu,
  Leaf,
  Mic,
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
