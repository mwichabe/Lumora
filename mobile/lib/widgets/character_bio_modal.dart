import 'package:flutter/material.dart';

import '../core/characters.dart';
import '../core/theme/colors.dart';
import '../core/theme/radii.dart';
import '../models/character.dart';
import 'avatar.dart';

void showCharacterBioModal(BuildContext context, CharacterWithFriendship character) {
  final info = characterInfo(character.name);
  showModalBottomSheet(
    context: context,
    backgroundColor: Colors.white,
    shape: const RoundedRectangleBorder(borderRadius: BorderRadius.vertical(top: Radius.circular(24))),
    builder: (context) => SafeArea(
      child: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            CircleAvatar(radius: 40, backgroundColor: colorFromHex(info.color), child: Text(character.emoji, style: const TextStyle(fontSize: 32))),
            const SizedBox(height: 12),
            Text(character.name, style: const TextStyle(fontSize: 20, fontWeight: FontWeight.w800)),
            Text(character.role, style: const TextStyle(color: LumoraColors.slatey)),
            const SizedBox(height: 12),
            if (info.quote.isNotEmpty)
              Container(
                padding: const EdgeInsets.all(14),
                decoration: BoxDecoration(color: LumoraColors.purpleLight, borderRadius: BorderRadius.circular(LumoraRadii.lg)),
                child: Text('"${info.quote}"', textAlign: TextAlign.center, style: const TextStyle(fontStyle: FontStyle.italic, color: LumoraColors.purple)),
              ),
            const SizedBox(height: 12),
            Text(character.personality, textAlign: TextAlign.center, style: const TextStyle(color: LumoraColors.slatey, fontSize: 13)),
            const SizedBox(height: 16),
            Row(mainAxisAlignment: MainAxisAlignment.center, children: [
              const Text('Friendship: ', style: TextStyle(fontWeight: FontWeight.w700)),
              for (var i = 0; i < 3; i++)
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 2),
                  child: Icon(Icons.star, size: 18, color: i < character.friendshipLevel ? LumoraColors.amber : LumoraColors.gray300),
                ),
            ]),
          ],
        ),
      ),
    ),
  );
}
