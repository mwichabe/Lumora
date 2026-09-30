import '../core/json_utils.dart';

class User {
  final int id;
  final String email;
  final String name;
  final String avatarColor;
  final String avatarUrl;
  final String targetLanguage;
  final String nativeLanguage;
  final String cefrLevel;
  final String levelName;
  final int dailyGoalXp;
  final int xp;
  final int xpToday;
  final int gems;
  final int hearts;
  final int streak;
  final int fluencyScore;
  final String league;
  final bool examUnlocked;

  const User({
    required this.id,
    required this.email,
    required this.name,
    required this.avatarColor,
    required this.avatarUrl,
    required this.targetLanguage,
    required this.nativeLanguage,
    required this.cefrLevel,
    required this.levelName,
    required this.dailyGoalXp,
    required this.xp,
    required this.xpToday,
    required this.gems,
    required this.hearts,
    required this.streak,
    required this.fluencyScore,
    required this.league,
    required this.examUnlocked,
  });

  factory User.fromJson(Map<String, dynamic> j) => User(
        id: asInt(j['id']),
        email: asString(j['email']),
        name: asString(j['name']),
        avatarColor: asString(j['avatarColor'], '#6C3FC5'),
        avatarUrl: asString(j['avatarUrl']),
        targetLanguage: asString(j['targetLanguage']),
        nativeLanguage: asString(j['nativeLanguage']),
        cefrLevel: asString(j['cefrLevel']),
        levelName: asString(j['levelName']),
        dailyGoalXp: asInt(j['dailyGoalXp']),
        xp: asInt(j['xp']),
        xpToday: asInt(j['xpToday']),
        gems: asInt(j['gems']),
        hearts: asInt(j['hearts']),
        streak: asInt(j['streak']),
        fluencyScore: asInt(j['fluencyScore']),
        league: asString(j['league']),
        examUnlocked: asBool(j['examUnlocked']),
      );

  Map<String, dynamic> toJson() => {
        'id': id,
        'email': email,
        'name': name,
        'avatarColor': avatarColor,
        'avatarUrl': avatarUrl,
        'targetLanguage': targetLanguage,
        'nativeLanguage': nativeLanguage,
        'cefrLevel': cefrLevel,
        'levelName': levelName,
        'dailyGoalXp': dailyGoalXp,
        'xp': xp,
        'xpToday': xpToday,
        'gems': gems,
        'hearts': hearts,
        'streak': streak,
        'fluencyScore': fluencyScore,
        'league': league,
        'examUnlocked': examUnlocked,
      };
}

class PaymentStatus {
  final bool paymentsEnabled;
  final String currency;
  final Map<String, int> prices;
  final Map<String, double> pricesUsd;
  final Map<String, bool> paid;

  const PaymentStatus({
    required this.paymentsEnabled,
    required this.currency,
    required this.prices,
    required this.pricesUsd,
    required this.paid,
  });

  factory PaymentStatus.fromJson(Map<String, dynamic> j) => PaymentStatus(
        paymentsEnabled: asBool(j['paymentsEnabled']),
        currency: asString(j['currency'], 'KES'),
        prices: asIntMap(j['prices']),
        pricesUsd: asDoubleMap(j['pricesUsd']),
        paid: asMap(j['paid']).map((k, v) => MapEntry(k, asBool(v))),
      );
}

class HeartsStatus {
  final int hearts;
  final int max;
  final bool full;
  final int secondsToNext;
  final int regenMinutes;
  final bool paymentsEnabled;
  final int refillPriceKes;
  final double refillPriceUsd;

  const HeartsStatus({
    required this.hearts,
    required this.max,
    required this.full,
    required this.secondsToNext,
    required this.regenMinutes,
    required this.paymentsEnabled,
    required this.refillPriceKes,
    required this.refillPriceUsd,
  });

  factory HeartsStatus.fromJson(Map<String, dynamic> j) => HeartsStatus(
        hearts: asInt(j['hearts']),
        max: asInt(j['max'], 5),
        full: asBool(j['full']),
        secondsToNext: asInt(j['secondsToNext']),
        regenMinutes: asInt(j['regenMinutes'], 30),
        paymentsEnabled: asBool(j['paymentsEnabled']),
        refillPriceKes: asInt(j['refillPriceKes']),
        refillPriceUsd: asDouble(j['refillPriceUsd']),
      );
}
