import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'colors.dart';
import 'radii.dart';

class AppTheme {
  AppTheme._();

  static ThemeData light() {
    final base = ThemeData.light(useMaterial3: true);
    return base.copyWith(
      scaffoldBackgroundColor: LumoraColors.cream,
      colorScheme: base.colorScheme.copyWith(
        primary: LumoraColors.purple,
        secondary: LumoraColors.amber,
        error: LumoraColors.coral,
        surface: Colors.white,
      ),
      textTheme: GoogleFonts.nunitoTextTheme(base.textTheme).apply(
        bodyColor: LumoraColors.ink,
        displayColor: LumoraColors.ink,
      ),
      appBarTheme: const AppBarTheme(
        backgroundColor: LumoraColors.purple,
        foregroundColor: Colors.white,
        elevation: 0,
        centerTitle: false,
      ),
      elevatedButtonTheme: ElevatedButtonThemeData(
        style: ElevatedButton.styleFrom(
          backgroundColor: LumoraColors.purple,
          foregroundColor: Colors.white,
          minimumSize: const Size.fromHeight(52),
          shape: const StadiumBorder(),
          textStyle: GoogleFonts.nunito(fontWeight: FontWeight.w800, fontSize: 16),
          elevation: 0,
        ),
      ),
      outlinedButtonTheme: OutlinedButtonThemeData(
        style: OutlinedButton.styleFrom(
          foregroundColor: LumoraColors.purple,
          side: const BorderSide(color: LumoraColors.purple, width: 2),
          minimumSize: const Size.fromHeight(52),
          shape: const StadiumBorder(),
          textStyle: GoogleFonts.nunito(fontWeight: FontWeight.w800, fontSize: 16),
        ),
      ),
      textButtonTheme: TextButtonThemeData(
        style: TextButton.styleFrom(
          foregroundColor: LumoraColors.purple,
          textStyle: GoogleFonts.nunito(fontWeight: FontWeight.w800, fontSize: 14),
        ),
      ),
      inputDecorationTheme: InputDecorationTheme(
        filled: true,
        fillColor: LumoraColors.gray50,
        contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(LumoraRadii.md),
          borderSide: BorderSide.none,
        ),
        enabledBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(LumoraRadii.md),
          borderSide: BorderSide.none,
        ),
        focusedBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(LumoraRadii.md),
          borderSide: const BorderSide(color: LumoraColors.purple, width: 2),
        ),
        errorBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(LumoraRadii.md),
          borderSide: const BorderSide(color: LumoraColors.coral, width: 1.5),
        ),
        hintStyle: GoogleFonts.nunito(color: LumoraColors.gray500),
      ),
      cardTheme: CardThemeData(
        color: Colors.white,
        elevation: 0,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(LumoraRadii.lg)),
      ),
      dividerTheme: const DividerThemeData(color: LumoraColors.gray100, thickness: 1),
      splashFactory: NoSplash.splashFactory,
      highlightColor: Colors.transparent,
    );
  }
}
