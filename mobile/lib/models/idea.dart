import '../core/json_utils.dart';
import 'chat.dart';

enum IdeaStatus { draft, underReview, approved, inProgress, completed, archived }

IdeaStatus ideaStatusFromString(String s) {
  switch (s) {
    case 'draft':
      return IdeaStatus.draft;
    case 'under_review':
      return IdeaStatus.underReview;
    case 'approved':
      return IdeaStatus.approved;
    case 'in_progress':
      return IdeaStatus.inProgress;
    case 'completed':
      return IdeaStatus.completed;
    case 'archived':
      return IdeaStatus.archived;
    default:
      return IdeaStatus.draft;
  }
}

String ideaStatusToString(IdeaStatus s) {
  switch (s) {
    case IdeaStatus.draft:
      return 'draft';
    case IdeaStatus.underReview:
      return 'under_review';
    case IdeaStatus.approved:
      return 'approved';
    case IdeaStatus.inProgress:
      return 'in_progress';
    case IdeaStatus.completed:
      return 'completed';
    case IdeaStatus.archived:
      return 'archived';
  }
}

String ideaStatusLabel(IdeaStatus s) {
  switch (s) {
    case IdeaStatus.draft:
      return 'Draft';
    case IdeaStatus.underReview:
      return 'Under Review';
    case IdeaStatus.approved:
      return 'Approved';
    case IdeaStatus.inProgress:
      return 'In Progress';
    case IdeaStatus.completed:
      return 'Completed';
    case IdeaStatus.archived:
      return 'Archived';
  }
}

class Idea {
  final int id;
  final String title;
  final String description;
  final IdeaStatus status;
  final ChatUser owner;
  final int upvotes;
  final int downvotes;
  final int score;
  final int myVote;
  final bool starred;
  final List<String> tags;
  final int messageCount;
  final String createdAt;
  final String lastActivity;
  final bool archived;
  final String archiveReason;
  final int? mergedIntoId;
  final int heat;

  const Idea({
    required this.id,
    required this.title,
    required this.description,
    required this.status,
    required this.owner,
    required this.upvotes,
    required this.downvotes,
    required this.score,
    required this.myVote,
    required this.starred,
    required this.tags,
    required this.messageCount,
    required this.createdAt,
    required this.lastActivity,
    required this.archived,
    required this.archiveReason,
    required this.mergedIntoId,
    required this.heat,
  });

  factory Idea.fromJson(Map<String, dynamic> j) => Idea(
        id: asInt(j['id']),
        title: asString(j['title']),
        description: asString(j['description']),
        status: ideaStatusFromString(asString(j['status'])),
        owner: ChatUser.fromJson(asMap(j['owner'])),
        upvotes: asInt(j['upvotes']),
        downvotes: asInt(j['downvotes']),
        score: asInt(j['score']),
        myVote: asInt(j['myVote']),
        starred: asBool(j['starred']),
        tags: asStringList(j['tags']),
        messageCount: asInt(j['messageCount']),
        createdAt: asString(j['createdAt']),
        lastActivity: asString(j['lastActivity']),
        archived: asBool(j['archived']),
        archiveReason: asString(j['archiveReason']),
        mergedIntoId: j['mergedIntoId'] == null ? null : asInt(j['mergedIntoId']),
        heat: asInt(j['heat']),
      );
}

class IdeaTagCount {
  final String tag;
  final int count;
  const IdeaTagCount({required this.tag, required this.count});
  factory IdeaTagCount.fromJson(Map<String, dynamic> j) =>
      IdeaTagCount(tag: asString(j['tag']), count: asInt(j['count']));
}

class IdeaBoard {
  final List<Idea> ideas;
  final List<IdeaTagCount> tags;
  final Map<String, int> counts;
  final int openIdeas;
  final int maxOpenIdeas;
  final bool crowded;

  const IdeaBoard({
    required this.ideas,
    required this.tags,
    required this.counts,
    required this.openIdeas,
    required this.maxOpenIdeas,
    required this.crowded,
  });

  factory IdeaBoard.fromJson(Map<String, dynamic> j) => IdeaBoard(
        ideas: asList(j['ideas'], (e) => Idea.fromJson(asMap(e))),
        tags: asList(j['tags'], (e) => IdeaTagCount.fromJson(asMap(e))),
        counts: asIntMap(j['counts']),
        openIdeas: asInt(j['openIdeas']),
        maxOpenIdeas: asInt(j['maxOpenIdeas']),
        crowded: asBool(j['crowded']),
      );
}

class IdeaEvent {
  final int id;
  final String kind;
  final String field;
  final String from;
  final String to;
  final String note;
  final ChatUser actor;
  final String at;

  const IdeaEvent({
    required this.id,
    required this.kind,
    required this.field,
    required this.from,
    required this.to,
    required this.note,
    required this.actor,
    required this.at,
  });

  factory IdeaEvent.fromJson(Map<String, dynamic> j) => IdeaEvent(
        id: asInt(j['id']),
        kind: asString(j['kind']),
        field: asString(j['field']),
        from: asString(j['from']),
        to: asString(j['to']),
        note: asString(j['note']),
        actor: ChatUser.fromJson(asMap(j['actor'])),
        at: asString(j['at']),
      );
}

class IdeaTask {
  final int id;
  final int ideaId;
  final String title;
  final String status; // todo | doing | done
  final String sprint;
  final String createdAt;
  final String? completedAt;

  const IdeaTask({
    required this.id,
    required this.ideaId,
    required this.title,
    required this.status,
    required this.sprint,
    required this.createdAt,
    required this.completedAt,
  });

  factory IdeaTask.fromJson(Map<String, dynamic> j) => IdeaTask(
        id: asInt(j['id']),
        ideaId: asInt(j['ideaId']),
        title: asString(j['title']),
        status: asString(j['status'], 'todo'),
        sprint: asString(j['sprint']),
        createdAt: asString(j['createdAt']),
        completedAt: j['completedAt'] == null ? null : asString(j['completedAt']),
      );
}

class SimilarIdea {
  final int id;
  final String title;
  final IdeaStatus status;
  final int score;
  final int messageCount;
  final double similarity;

  const SimilarIdea({
    required this.id,
    required this.title,
    required this.status,
    required this.score,
    required this.messageCount,
    required this.similarity,
  });

  factory SimilarIdea.fromJson(Map<String, dynamic> j) => SimilarIdea(
        id: asInt(j['id']),
        title: asString(j['title']),
        status: ideaStatusFromString(asString(j['status'])),
        score: asInt(j['score']),
        messageCount: asInt(j['messageCount']),
        similarity: asDouble(j['similarity']),
      );
}

/// One move available from the idea's current status, as the server sees it
/// for this viewer. Disallowed moves still come back, with [reason], so the UI
/// can explain what has to happen next.
class IdeaTransition {
  final IdeaStatus to;
  final String label;
  final String hint;
  final bool allowed;
  final String reason;
  final bool primary;

  const IdeaTransition({
    required this.to,
    required this.label,
    required this.hint,
    required this.allowed,
    required this.reason,
    required this.primary,
  });

  factory IdeaTransition.fromJson(Map<String, dynamic> j) => IdeaTransition(
        to: ideaStatusFromString(asString(j['to'])),
        label: asString(j['label']),
        hint: asString(j['hint']),
        allowed: asBool(j['allowed']),
        reason: asString(j['reason']),
        primary: asBool(j['primary']),
      );
}

class IdeaDetail {
  final Idea idea;
  final List<IdeaEvent> history;
  final List<IdeaTask> tasks;
  final List<ChatUser> participants;
  final List<Idea> mergedIn;
  final bool canEdit;
  final List<IdeaStatus> statusFlow;
  final List<SimilarIdea> similar;
  final List<IdeaTransition> transitions;
  final String nextStep;
  final int openTasks;

  const IdeaDetail({
    required this.idea,
    required this.history,
    required this.tasks,
    required this.participants,
    required this.mergedIn,
    required this.canEdit,
    required this.statusFlow,
    required this.similar,
    required this.transitions,
    required this.nextStep,
    required this.openTasks,
  });

  factory IdeaDetail.fromJson(Map<String, dynamic> j) => IdeaDetail(
        idea: Idea.fromJson(asMap(j['idea'])),
        history: asList(j['history'], (e) => IdeaEvent.fromJson(asMap(e))),
        tasks: asList(j['tasks'], (e) => IdeaTask.fromJson(asMap(e))),
        participants: asList(j['participants'], (e) => ChatUser.fromJson(asMap(e))),
        mergedIn: asList(j['mergedIn'], (e) => Idea.fromJson(asMap(e))),
        canEdit: asBool(j['canEdit']),
        statusFlow: asList(j['statusFlow'], (e) => ideaStatusFromString(asString(e))),
        similar: asList(j['similar'], (e) => SimilarIdea.fromJson(asMap(e))),
        transitions: asList(j['transitions'], (e) => IdeaTransition.fromJson(asMap(e))),
        nextStep: asString(j['nextStep']),
        openTasks: asInt(j['openTasks']),
      );
}

class IdeaReaction {
  final String emoji;
  final int count;
  final bool mine;
  const IdeaReaction({required this.emoji, required this.count, required this.mine});
  factory IdeaReaction.fromJson(Map<String, dynamic> j) => IdeaReaction(
        emoji: asString(j['emoji']),
        count: asInt(j['count']),
        mine: asBool(j['mine']),
      );
}

class IdeaMessage {
  final int id;
  final int ideaId;
  final int? parentId;
  final ChatUser? author;
  final String kind; // text|image|voice|code
  final String body;
  final String fileName;
  final String url;
  final int width;
  final int height;
  final int duration;
  final bool anonymous;
  final MessageTranslation? translation;
  final List<IdeaReaction> reactions;
  final List<IdeaMessage> replies;
  final int replyCount;
  final bool mine;
  final bool edited;
  final bool deleted;
  final bool canEdit;
  final String createdAt;

  const IdeaMessage({
    required this.id,
    required this.ideaId,
    required this.parentId,
    required this.author,
    required this.kind,
    required this.body,
    required this.fileName,
    required this.url,
    required this.width,
    required this.height,
    required this.duration,
    required this.anonymous,
    required this.translation,
    required this.reactions,
    required this.replies,
    required this.replyCount,
    required this.mine,
    required this.edited,
    required this.deleted,
    required this.canEdit,
    required this.createdAt,
  });

  factory IdeaMessage.fromJson(Map<String, dynamic> j) => IdeaMessage(
        id: asInt(j['id']),
        ideaId: asInt(j['ideaId']),
        parentId: j['parentId'] == null ? null : asInt(j['parentId']),
        author: j['author'] == null ? null : ChatUser.fromJson(asMap(j['author'])),
        kind: asString(j['kind'], 'text'),
        body: asString(j['body']),
        fileName: asString(j['fileName']),
        url: asString(j['url']),
        width: asInt(j['width']),
        height: asInt(j['height']),
        duration: asInt(j['duration']),
        anonymous: asBool(j['anonymous']),
        translation:
            j['translation'] == null ? null : MessageTranslation.fromJson(asMap(j['translation'])),
        reactions: asList(j['reactions'], (e) => IdeaReaction.fromJson(asMap(e))),
        replies: asList(j['replies'], (e) => IdeaMessage.fromJson(asMap(e))),
        replyCount: asInt(j['replyCount']),
        mine: asBool(j['mine']),
        edited: asBool(j['edited']),
        deleted: asBool(j['deleted']),
        canEdit: asBool(j['canEdit']),
        createdAt: asString(j['createdAt']),
      );
}

class BrainstormSession {
  final int id;
  final String topic;
  final String endsAt;
  final int secondsRemaining;

  const BrainstormSession({
    required this.id,
    required this.topic,
    required this.endsAt,
    required this.secondsRemaining,
  });

  factory BrainstormSession.fromJson(Map<String, dynamic> j) => BrainstormSession(
        id: asInt(j['id']),
        topic: asString(j['topic']),
        endsAt: asString(j['endsAt']),
        secondsRemaining: asInt(j['secondsRemaining']),
      );
}

class IdeaThread {
  final List<IdeaMessage> messages;
  final Idea idea;
  final BrainstormSession? brainstorm;
  final List<String> reactions;

  const IdeaThread({
    required this.messages,
    required this.idea,
    required this.brainstorm,
    required this.reactions,
  });

  factory IdeaThread.fromJson(Map<String, dynamic> j) => IdeaThread(
        messages: asList(j['messages'], (e) => IdeaMessage.fromJson(asMap(e))),
        idea: Idea.fromJson(asMap(j['idea'])),
        brainstorm:
            j['brainstorm'] == null ? null : BrainstormSession.fromJson(asMap(j['brainstorm'])),
        reactions: asStringList(j['reactions']),
      );
}

class ThreadSummaryPoint {
  final String text;
  final ChatUser author;
  final String at;
  const ThreadSummaryPoint({required this.text, required this.author, required this.at});
  factory ThreadSummaryPoint.fromJson(Map<String, dynamic> j) => ThreadSummaryPoint(
        text: asString(j['text']),
        author: ChatUser.fromJson(asMap(j['author'])),
        at: asString(j['at']),
      );
}

class ThreadSummaryQuestion {
  final String text;
  final int authorId;
  final String at;
  const ThreadSummaryQuestion({required this.text, required this.authorId, required this.at});
  factory ThreadSummaryQuestion.fromJson(Map<String, dynamic> j) => ThreadSummaryQuestion(
        text: asString(j['text']),
        authorId: asInt(j['authorId']),
        at: asString(j['at']),
      );
}

class ThreadSummary {
  final String gist;
  final List<ThreadSummaryPoint> keyPoints;
  final List<ThreadSummaryQuestion> questions;
  final int messageCount;
  final bool generated;

  const ThreadSummary({
    required this.gist,
    required this.keyPoints,
    required this.questions,
    required this.messageCount,
    required this.generated,
  });

  factory ThreadSummary.fromJson(Map<String, dynamic> j) => ThreadSummary(
        gist: asString(j['gist']),
        keyPoints: asList(j['keyPoints'], (e) => ThreadSummaryPoint.fromJson(asMap(e))),
        questions: asList(j['questions'], (e) => ThreadSummaryQuestion.fromJson(asMap(e))),
        messageCount: asInt(j['messageCount']),
        generated: asBool(j['generated']),
      );
}
