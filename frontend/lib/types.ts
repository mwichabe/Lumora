// Types mirror the JSON returned by the Lumora Go API.

export interface User {
  id: number;
  email: string;
  name: string;
  avatarColor: string;
  avatarUrl: string;
  targetLanguage: string;
  nativeLanguage: string;
  cefrLevel: string;
  levelName: string;
  dailyGoalXp: number;
  xp: number;
  xpToday: number;
  gems: number;
  hearts: number;
  streak: number;
  fluencyScore: number;
  league: string;
  examUnlocked?: boolean;
}

export interface PaymentStatus {
  paymentsEnabled: boolean;
  currency: string;
  prices: Record<string, number>; // level -> price in KES
  pricesUsd: Record<string, number>; // level -> approx price in USD
  paid: Record<string, boolean>; // level -> has an unconsumed paid attempt
}

export interface HeartsStatus {
  hearts: number;
  max: number;
  full: boolean;
  secondsToNext: number; // until the next heart regenerates (0 when full)
  regenMinutes: number;
  paymentsEnabled: boolean;
  refillPriceKes: number;
  refillPriceUsd: number;
}

export type ExerciseType =
  | "translate"
  | "multiple_choice"
  | "listen"
  | "match"
  | "speak"
  | "fill"
  | "write"
  | "character";

export interface Exercise {
  id: number;
  lessonId: number;
  type: ExerciseType;
  orderIndex: number;
  prompt: string;
  question: string;
  options: string[] | null;
  correctAnswer: string;
  character: string;
}

export interface VocabItem {
  id: number;
  lessonId: number;
  orderIndex: number;
  word: string;
  translation: string;
  example: string;
  exampleTranslation: string;
  speaker: string;
}

export interface Mistake {
  id: number;
  userId: number;
  language: string;
  prompt: string;
  question: string;
  correctAnswer: string;
}

export interface Certificate {
  id: number;
  userId: number;
  userName: string;
  language: string;
  level: string;
  score: number;
  listening: number;
  reading: number;
  writing: number;
  speaking: number;
  serial: string;
  issuedAt: string;
}

export interface ExamSectionWeights {
  listening: number;
  reading: number;
  writing: number;
  speaking: number;
}

export interface ExamResult {
  passed: boolean;
  alreadyTaken?: boolean;
  overall: number;
  level: string;
  passMark?: number;
  weights?: ExamSectionWeights;
  sections?: {
    listening: number;
    reading: number;
    writing: number;
    speaking: number;
  };
  certificate?: Certificate;
}

export interface ExamMeta {
  weights: ExamSectionWeights;
  levels: Record<string, { passMark: number; durationSeconds: number }>;
}

export interface PaperQuestion {
  question: string;
  options: string[];
  correctAnswer: string;
}

export interface PaperLine {
  character: string;
  text: string;
  translation: string;
}

export interface ExamPaper {
  ready: boolean;
  language: string;
  level: string;
  durationSeconds: number;
  passMark: number;
  weights: ExamSectionWeights;
  listening: {
    title: string;
    lines: PaperLine[];
    questions: PaperQuestion[];
  } | null;
  reading: {
    title: string;
    paragraphs: string[];
    questions: PaperQuestion[];
  } | null;
  writing: { prompt: string; minWords: number };
  speaking: { phrase: string; speaker: string; translation: string };
}

export interface CertVerification {
  valid: boolean;
  certificate?: {
    userName: string;
    language: string;
    level: string;
    score: number;
    serial: string;
    issuedAt: string;
  };
}

export interface AppNotification {
  id: number;
  kind: string;
  emoji: string;
  tint: string;
  title: string;
  body: string;
  link: string;
  read: boolean;
  createdAt: string;
}

export interface ChatUser {
  id: number;
  name: string;
  avatarColor: string;
  avatarUrl: string;
  levelName: string;
}

export interface ChatMessage {
  id: number;
  senderId: number;
  recipientId: number;
  body: string;
  read: boolean;
  createdAt: string;
}

export interface ChatThread {
  user: ChatUser;
  lastMessage: string;
  lastAt: string;
  unread: number;
}

export interface Lesson {
  id: number;
  skillId: number;
  title: string;
  orderIndex: number;
  xpReward: number;
  vocab?: VocabItem[];
  exercises?: Exercise[];
}

export interface ListeningLine {
  id: number;
  sessionId: number;
  orderIndex: number;
  character: string;
  text: string;
  translation: string;
}

export interface ListeningQuestion {
  id: number;
  sessionId: number;
  orderIndex: number;
  prompt: string;
  question: string;
  options: string[] | null;
  correctAnswer: string;
}

export interface ReadingLine {
  id: number;
  sessionId: number;
  orderIndex: number;
  text: string;
  translation: string;
}

export interface ReadingQuestion {
  id: number;
  sessionId: number;
  orderIndex: number;
  prompt: string;
  question: string;
  options: string[] | null;
  correctAnswer: string;
}

export interface ReadingSession {
  id: number;
  language: string;
  unit: string;
  title: string;
  description: string;
  orderIndex: number;
  xpReward: number;
  lines?: ReadingLine[];
  questions?: ReadingQuestion[];
}

export interface ListeningMatch {
  id: number;
  sessionId: number;
  orderIndex: number;
  word: string;
  translation: string;
}

export interface ListeningSession {
  id: number;
  language: string;
  unit: string;
  title: string;
  description: string;
  orderIndex: number;
  xpReward: number;
  matches?: ListeningMatch[];
  lines?: ListeningLine[];
  questions?: ListeningQuestion[];
}

export interface Skill {
  id: number;
  language: string;
  unit: string;
  title: string;
  description: string;
  icon: string;
  color: string;
  orderIndex: number;
  requiredXp: number;
  lessons?: Lesson[];
  unlocked: boolean;
  completed: boolean;
  lessonCount: number;
  completedCount: number;
}

export interface Quest {
  id: number;
  title: string;
  description: string;
  icon: string;
  xpReward: number;
  target: number;
}

export interface UserQuest {
  id: number;
  questId: number;
  date: string;
  progress: number;
  completed: boolean;
  quest?: Quest;
}

export interface CharacterWithFriendship {
  id: number;
  name: string;
  species: string;
  role: string;
  personality: string;
  color: string;
  emoji: string;
  friendshipLevel: number;
}

export interface LeaderRow {
  id: number;
  name: string;
  xp: number;
  streak: number;
  avatarColor: string;
  avatarUrl: string;
  language: string;
  isUser: boolean;
  rank: number;
}

export interface HomeData {
  user: User;
  nextLesson: Lesson | null;
  nextSkill: Skill | null;
  quests: UserQuest[];
}
