import '../core/json_utils.dart';

class ChatUser {
  final int id;
  final String name;
  final String avatarColor;
  final String avatarUrl;
  final String levelName;

  const ChatUser({
    required this.id,
    required this.name,
    required this.avatarColor,
    required this.avatarUrl,
    required this.levelName,
  });

  factory ChatUser.fromJson(Map<String, dynamic> j) => ChatUser(
        id: asInt(j['id']),
        name: asString(j['name']),
        avatarColor: asString(j['avatarColor'], '#6C3FC5'),
        avatarUrl: asString(j['avatarUrl']),
        levelName: asString(j['levelName']),
      );
}

class MessageTranslation {
  final String lang;
  final String langName;
  final String text;
  final bool pending;

  const MessageTranslation({
    required this.lang,
    required this.langName,
    required this.text,
    required this.pending,
  });

  factory MessageTranslation.fromJson(Map<String, dynamic> j) => MessageTranslation(
        lang: asString(j['lang']),
        langName: asString(j['langName']),
        text: asString(j['text']),
        pending: asBool(j['pending']),
      );
}

class ChatMessage {
  final int id;
  final int senderId;
  final int recipientId;
  final String kind; // text | image
  final String body;
  final String url;
  final String fileName;
  final int width;
  final int height;
  final bool read;
  final bool mine;
  final bool edited;
  final bool deleted;
  final bool canEdit;
  final String createdAt;
  final MessageTranslation? translation;

  const ChatMessage({
    required this.id,
    required this.senderId,
    required this.recipientId,
    required this.kind,
    required this.body,
    required this.url,
    required this.fileName,
    required this.width,
    required this.height,
    required this.read,
    required this.mine,
    required this.edited,
    required this.deleted,
    required this.canEdit,
    required this.createdAt,
    required this.translation,
  });

  factory ChatMessage.fromJson(Map<String, dynamic> j) => ChatMessage(
        id: asInt(j['id']),
        senderId: asInt(j['senderId']),
        recipientId: asInt(j['recipientId']),
        kind: asString(j['kind'], 'text'),
        body: asString(j['body']),
        url: asString(j['url']),
        fileName: asString(j['fileName']),
        width: asInt(j['width']),
        height: asInt(j['height']),
        read: asBool(j['read']),
        mine: asBool(j['mine']),
        edited: asBool(j['edited']),
        deleted: asBool(j['deleted']),
        canEdit: asBool(j['canEdit']),
        createdAt: asString(j['createdAt']),
        translation:
            j['translation'] == null ? null : MessageTranslation.fromJson(asMap(j['translation'])),
      );
}

class ChatThread {
  final ChatUser user;
  final String lastMessage;
  final String lastAt;
  final int unread;

  const ChatThread({
    required this.user,
    required this.lastMessage,
    required this.lastAt,
    required this.unread,
  });

  factory ChatThread.fromJson(Map<String, dynamic> j) => ChatThread(
        user: ChatUser.fromJson(asMap(j['user'])),
        lastMessage: asString(j['lastMessage']),
        lastAt: asString(j['lastAt']),
        unread: asInt(j['unread']),
      );
}
