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

/**
 * Attached to a message that wasn't written in English. Absent entirely on
 * English (or undetermined) messages, so the UI renders nothing at all in the
 * overwhelmingly common case.
 */
export interface MessageTranslation {
  lang: string;      // ISO 639-1 of the original
  langName: string;  // "Spanish"
  text: string;      // English — empty while still being translated
  pending: boolean;  // a translation is on its way
}

export interface ChatMessage {
  id: number;
  senderId: number;
  recipientId: number;
  kind: "text" | "image";
  body: string;
  /** Attachment endpoint, when the message carries one. */
  url: string;
  fileName: string;
  width: number;
  height: number;
  read: boolean;
  mine: boolean;
  edited: boolean;
  /** Soft-deleted: render a tombstone, not nothing. */
  deleted: boolean;
  canEdit: boolean;
  createdAt: string;
  translation?: MessageTranslation;
}

// --- the ideas workspace -----------------------------------------------------

export type IdeaStatus =
  | "draft"
  | "under_review"
  | "approved"
  | "in_progress"
  | "completed"
  | "archived";

export interface Idea {
  id: number;
  title: string;
  description: string;
  status: IdeaStatus;
  owner: ChatUser;
  upvotes: number;
  downvotes: number;
  score: number;
  /** The viewer's own vote: 1, -1 or 0. */
  myVote: number;
  starred: boolean;
  tags: string[];
  messageCount: number;
  createdAt: string;
  lastActivity: string;
  archived: boolean;
  archiveReason: string;
  mergedIntoId: number | null;
  heat: number;
}

export interface IdeaBoard {
  ideas: Idea[];
  tags: { tag: string; count: number }[];
  counts: Record<string, number>;
  openIdeas: number;
  maxOpenIdeas: number;
  /** True once the board has more open ideas than anyone will read. */
  crowded: boolean;
}

export interface IdeaEvent {
  id: number;
  kind: string;
  field: string;
  from: string;
  to: string;
  note: string;
  actor: ChatUser;
  at: string;
}

export interface IdeaTask {
  id: number;
  ideaId: number;
  title: string;
  status: "todo" | "doing" | "done";
  sprint: string;
  createdAt: string;
  completedAt: string | null;
}

export interface SimilarIdea {
  id: number;
  title: string;
  status: IdeaStatus;
  score: number;
  messageCount: number;
  similarity: number;
}

export interface IdeaDetail {
  idea: Idea;
  history: IdeaEvent[];
  tasks: IdeaTask[];
  participants: ChatUser[];
  mergedIn: Idea[];
  canEdit: boolean;
  statusFlow: IdeaStatus[];
  similar: SimilarIdea[];
}

export interface IdeaReaction {
  emoji: string;
  count: number;
  mine: boolean;
}

export interface IdeaMessage {
  id: number;
  ideaId: number;
  parentId: number | null;
  /** Null while a silent brainstorm is running — anonymity is enforced server-side. */
  author: ChatUser | null;
  kind: "text" | "image" | "voice" | "code";
  body: string;
  fileName: string;
  url: string;
  width: number;
  height: number;
  duration: number;
  anonymous: boolean;
  translation?: MessageTranslation;
  reactions: IdeaReaction[];
  replies: IdeaMessage[];
  replyCount: number;
  mine: boolean;
  edited: boolean;
  deleted: boolean;
  canEdit: boolean;
  createdAt: string;
}

export interface BrainstormSession {
  id: number;
  topic: string;
  endsAt: string;
  secondsRemaining: number;
}

export interface IdeaThread {
  messages: IdeaMessage[];
  idea: Idea;
  brainstorm: BrainstormSession | null;
  reactions: string[];
}

export interface ThreadSummary {
  gist: string;
  keyPoints: { text: string; author: ChatUser; at: string }[];
  questions: { text: string; authorId: number; at: string }[];
  messageCount: number;
  /** False when the thread was too thin to summarise — say so, don't invent one. */
  generated: boolean;
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

// --- the weekly league -------------------------------------------------------
//
// Ten tiers, pods of thirty, one week per season. `points` are NOT raw XP: the
// backend weights every activity by difficulty, accuracy, consistency and how
// early in the week it was earned. See backend/controllers/league_engine.go.

export interface LeagueTier {
  index: number;
  name: string;
  tint: string;
  promoteTop: number;
  demoteBottom: number;
  goldGems: number;
  groupBonus: number;
}

export type LeagueZone = "promote" | "hold" | "demote";

export interface LeagueRow {
  id: number;
  name: string;
  points: number;
  rawXp: number;
  streak: number;
  avatarColor: string;
  avatarUrl: string;
  language: string;
  rank: number;
  zone: LeagueZone;
  isUser: boolean;
  accuracy: number;
  perfectRuns: number;
  fairPlay: boolean;
  flagged: boolean;
  reported: boolean;
}

export interface LeagueStandings {
  seasonId: string;
  endsAt: string;
  secondsRemaining: number;
  tiers: LeagueTier[];
  tier: LeagueTier;
  podSize: number;
  casual: boolean;
  joined: boolean;
  rows: LeagueRow[];
  userRank?: number;
  promoteTop?: number;
  demoteBottom?: number;
  stage?: string;
  groupGoal?: {
    target: number;
    current: number;
    hit: boolean;
    bonus: number;
  };
  me?: {
    points: number;
    rawXp: number;
    activities: number;
    perfectRuns: number;
    flagged: boolean;
    flagReason: string;
  };
  you: {
    integrity: number;
    fairPlay: boolean;
    trophies: number;
    best: string;
    bestTier: number;
    casual: boolean;
  };
}

export type LeagueOutcome =
  | "promoted"
  | "held"
  | "demoted"
  | "qualified"
  | "advanced"
  | "champion"
  | "eliminated";

/** One settled season, replayed once as the end-of-race ceremony. */
export interface LeagueResult {
  seasonId: string;
  result: LeagueOutcome;
  rank: number;
  podSize: number;
  points: number;
  rawXp: number;
  activities: number;
  accuracy: number;
  perfectRuns: number;
  gems: number;
  groupGoalHit: boolean;
  flagged: boolean;
  flagReason: string;
  stage: string;
  from: LeagueTier;
  to: LeagueTier;
  podium: {
    rank: number;
    name: string;
    points: number;
    avatarColor: string;
    avatarUrl: string;
    isUser: boolean;
  }[];
  trophies: number;
}

export interface LeagueHistoryEntry {
  seasonId: string;
  tier: LeagueTier;
  to: LeagueTier;
  rank: number;
  points: number;
  result: LeagueOutcome;
  gems: number;
  stage: string;
}

export interface HomeData {
  user: User;
  nextLesson: Lesson | null;
  nextSkill: Skill | null;
  quests: UserQuest[];
}
