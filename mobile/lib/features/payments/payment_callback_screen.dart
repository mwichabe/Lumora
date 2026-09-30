import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../core/network/api_client.dart';
import '../../core/theme/colors.dart';
import '../../providers/auth_provider.dart';
import '../../widgets/fox_mascot.dart';
import '../../widgets/lumora_button.dart';

enum _State { checking, success, failed }

/// Reached if a Paystack redirect deep-links back into the app with
/// ?reference=… (frontend/app/payment/callback). Payments normally finish
/// inside the in-app checkout (checkout_screen.dart) and never get here, but
/// this screen still verifies and reacts correctly when the redirect does
/// land here.
class PaymentCallbackScreen extends ConsumerStatefulWidget {
  final String? reference;
  const PaymentCallbackScreen({super.key, this.reference});

  @override
  ConsumerState<PaymentCallbackScreen> createState() => _PaymentCallbackScreenState();
}

class _PaymentCallbackScreenState extends ConsumerState<PaymentCallbackScreen> {
  _State _state = _State.checking;
  String _product = '';

  @override
  void initState() {
    super.initState();
    _verify();
  }

  Future<void> _verify() async {
    final ref_ = widget.reference;
    if (ref_ == null || ref_.isEmpty) {
      setState(() => _state = _State.failed);
      return;
    }
    try {
      final r = await ApiClient.instance.verifyPayment(ref_);
      setState(() {
        _product = r.product;
        _state = r.success ? _State.success : _State.failed;
      });
      ref.read(authProvider.notifier).refresh();
    } catch (_) {
      setState(() => _state = _State.failed);
    }
  }

  @override
  Widget build(BuildContext context) {
    final isHearts = _product == 'hearts_refill';
    return Scaffold(
      backgroundColor: LumoraColors.cream,
      body: Center(
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: switch (_state) {
              _State.checking => const [
                  FoxMascot(size: 110, glow: true),
                  SizedBox(height: 20),
                  Text('Confirming your payment…', style: TextStyle(fontSize: 18, fontWeight: FontWeight.w800)),
                  SizedBox(height: 4),
                  Text('Hang tight, this only takes a moment.', style: TextStyle(color: LumoraColors.slatey)),
                ],
              _State.success => [
                  Container(width: 80, height: 80, decoration: BoxDecoration(color: LumoraColors.teal.withValues(alpha: 0.1), shape: BoxShape.circle),
                      child: const Icon(Icons.check_circle, color: LumoraColors.teal, size: 44)),
                  const SizedBox(height: 16),
                  Text(isHearts ? 'Hearts refilled! ❤️' : "You're all set! 🎉", style: const TextStyle(fontSize: 22, fontWeight: FontWeight.w800)),
                  const SizedBox(height: 4),
                  Text(isHearts ? 'Your hearts are full again — jump back into your lesson.' : 'Your exam attempt is ready — take it on the web.',
                      textAlign: TextAlign.center, style: const TextStyle(color: LumoraColors.slatey)),
                  const SizedBox(height: 20),
                  LumoraButton(label: isHearts ? 'Continue learning' : 'How to take it', full: true, onPressed: () => context.go(isHearts ? '/learn' : '/exam')),
                  const SizedBox(height: 8),
                  LumoraButton(label: 'Back to home', full: true, variant: LumoraButtonVariant.outline, onPressed: () => context.go('/home')),
                ],
              _State.failed => [
                  Container(width: 80, height: 80, decoration: BoxDecoration(color: LumoraColors.coral.withValues(alpha: 0.1), shape: BoxShape.circle),
                      child: const Icon(Icons.cancel, color: LumoraColors.coral, size: 44)),
                  const SizedBox(height: 16),
                  const Text('Payment not completed', style: TextStyle(fontSize: 22, fontWeight: FontWeight.w800)),
                  const SizedBox(height: 4),
                  const Text("We couldn't confirm your payment. If you were charged, it will be applied automatically — otherwise you can try again.",
                      textAlign: TextAlign.center, style: TextStyle(color: LumoraColors.slatey)),
                  const SizedBox(height: 20),
                  LumoraButton(label: 'Try again', full: true, onPressed: () => context.go('/exam')),
                  const SizedBox(height: 8),
                  LumoraButton(label: 'Back to profile', full: true, variant: LumoraButtonVariant.outline, onPressed: () => context.go('/profile')),
                ],
            },
          ),
        ),
      ),
    );
  }
}
