import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:url_launcher/url_launcher.dart';
import 'package:webview_flutter/webview_flutter.dart';

import '../../core/network/api_client.dart';
import '../../core/network/api_exception.dart';
import '../../core/theme/colors.dart';
import '../../widgets/fox_mascot.dart';
import '../../widgets/lumora_button.dart';

/// Runs a Paystack checkout without leaving the app and resolves to `true`
/// once the backend has confirmed the payment.
///
/// Pass [level] for an exam attempt or `product: 'hearts'` for a hearts
/// refill. A `false` result means "not confirmed here" — the user backed out,
/// or the payment is still settling. The Paystack webhook fulfils a late
/// payment server-side, so callers should re-read their own state (payment
/// status, hearts) after this returns rather than trust `false` as final.
Future<bool> payWithPaystack(BuildContext context, {String? level, String? product}) async {
  final paid = await Navigator.of(context, rootNavigator: true).push<bool>(
    MaterialPageRoute(
      fullscreenDialog: true,
      builder: (_) => _CheckoutScreen(level: level, product: product),
    ),
  );
  return paid ?? false;
}

/// webview_flutter only ships Android and iOS implementations; everywhere else
/// (desktop, web) the checkout opens in the system browser instead.
bool get _inAppCheckout =>
    !kIsWeb && (defaultTargetPlatform == TargetPlatform.android || defaultTargetPlatform == TargetPlatform.iOS);

enum _Stage { starting, paying, browser, confirming, failed }

class _CheckoutScreen extends StatefulWidget {
  final String? level;
  final String? product;
  const _CheckoutScreen({this.level, this.product});

  @override
  State<_CheckoutScreen> createState() => _CheckoutScreenState();
}

class _CheckoutScreenState extends State<_CheckoutScreen> {
  _Stage _stage = _Stage.starting;
  String _failure = '';
  String? _reference;
  WebViewController? _web;
  int _progress = 0;

  @override
  void initState() {
    super.initState();
    _start();
  }

  Future<void> _start() async {
    setState(() {
      _stage = _Stage.starting;
      _progress = 0;
    });
    try {
      final (url, reference) =
          await ApiClient.instance.initializePayment(level: widget.level, product: widget.product);
      if (!mounted) return;
      if (url == null || reference == null) {
        _fail("We couldn't start the payment. Please try again.");
        return;
      }
      _reference = reference;

      if (!_inAppCheckout) {
        await launchUrl(Uri.parse(url), mode: LaunchMode.externalApplication);
        if (mounted) setState(() => _stage = _Stage.browser);
        return;
      }

      _web = WebViewController()
        ..setJavaScriptMode(JavaScriptMode.unrestricted)
        ..setNavigationDelegate(NavigationDelegate(
          onProgress: (p) {
            if (mounted) setState(() => _progress = p);
          },
          // Paystack finishes by sending the browser to the web app's
          // /payment/callback page. There's nothing to see there inside the
          // app, so stop the navigation and confirm the payment natively.
          onNavigationRequest: (request) {
            if (_isReturnUrl(request.url)) {
              _confirm();
              return NavigationDecision.prevent;
            }
            return NavigationDecision.navigate;
          },
          // Server-side redirects don't always reach onNavigationRequest on
          // Android, so catch the return URL here as well.
          onPageStarted: (url) {
            if (_isReturnUrl(url)) _confirm();
          },
        ))
        ..loadRequest(Uri.parse(url));
      setState(() => _stage = _Stage.paying);
    } on ApiException catch (e) {
      _fail(e.status == 503 ? 'Payments are not available right now.' : "We couldn't start the payment. Please try again.");
    } catch (_) {
      _fail("We couldn't start the payment. Check your connection and try again.");
    }
  }

  bool _isReturnUrl(String url) {
    final uri = Uri.tryParse(url);
    return uri != null && uri.path.endsWith('/payment/callback');
  }

  void _fail(String message) {
    if (!mounted) return;
    setState(() {
      _failure = message;
      _stage = _Stage.failed;
    });
  }

  /// Asks the backend to verify the reference with Paystack. Mobile-money
  /// payments can take a moment to settle, so a "not yet" gets a few retries
  /// before it's reported as not completed.
  Future<void> _confirm() async {
    final reference = _reference;
    if (reference == null || _stage == _Stage.confirming) return;
    setState(() => _stage = _Stage.confirming);

    for (var attempt = 0; attempt < 4; attempt++) {
      if (attempt > 0) await Future<void>.delayed(const Duration(seconds: 2));
      if (!mounted) return;
      try {
        final result = await ApiClient.instance.verifyPayment(reference);
        if (!mounted) return;
        if (result.success) {
          Navigator.of(context).pop(true);
          return;
        }
      } catch (_) {
        // network blip — fall through and retry
      }
    }
    _fail("We couldn't confirm your payment. If you were charged, it will be applied automatically within a few minutes.");
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: LumoraColors.cream,
      appBar: AppBar(
        backgroundColor: Colors.white,
        surfaceTintColor: Colors.white,
        elevation: 0,
        leading: IconButton(
          icon: const Icon(Icons.close, color: LumoraColors.ink),
          tooltip: 'Close',
          onPressed: () => Navigator.of(context).pop(false),
        ),
        title: const Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(Icons.lock_rounded, size: 16, color: LumoraColors.teal),
            SizedBox(width: 6),
            Text('Secure checkout', style: TextStyle(fontSize: 16, fontWeight: FontWeight.w800, color: LumoraColors.ink)),
          ],
        ),
        centerTitle: true,
        bottom: PreferredSize(
          preferredSize: const Size.fromHeight(3),
          child: _stage == _Stage.paying && _progress < 100
              ? LinearProgressIndicator(
                  value: _progress / 100,
                  minHeight: 3,
                  backgroundColor: LumoraColors.purpleLight,
                  valueColor: const AlwaysStoppedAnimation(LumoraColors.purple),
                )
              : const SizedBox(height: 3),
        ),
      ),
      body: SafeArea(child: _body()),
    );
  }

  Widget _body() {
    switch (_stage) {
      case _Stage.paying:
        return WebViewWidget(controller: _web!);
      case _Stage.starting:
        return const _Message(title: 'Opening secure checkout…', detail: 'This only takes a moment.');
      case _Stage.confirming:
        return const _Message(title: 'Confirming your payment…', detail: "Hang tight — please don't close this screen.");
      case _Stage.browser:
        return _Message(
          title: 'Finish paying in your browser',
          detail: "We've opened the checkout in your browser. Come back here when you're done.",
          actions: [
            LumoraButton(label: "I've paid — continue", full: true, onPressed: _confirm),
            const SizedBox(height: 8),
            LumoraButton(
              label: 'Cancel',
              full: true,
              variant: LumoraButtonVariant.outline,
              onPressed: () => Navigator.of(context).pop(false),
            ),
          ],
        );
      case _Stage.failed:
        return _Message(
          title: 'Payment not completed',
          detail: _failure,
          failed: true,
          actions: [
            LumoraButton(label: 'Try again', full: true, onPressed: _start),
            const SizedBox(height: 8),
            LumoraButton(
              label: 'Close',
              full: true,
              variant: LumoraButtonVariant.outline,
              onPressed: () => Navigator.of(context).pop(false),
            ),
          ],
        );
    }
  }
}

class _Message extends StatelessWidget {
  final String title;
  final String detail;
  final bool failed;
  final List<Widget> actions;
  const _Message({required this.title, required this.detail, this.failed = false, this.actions = const []});

  @override
  Widget build(BuildContext context) {
    return Center(
      child: SingleChildScrollView(
        padding: const EdgeInsets.all(24),
        child: ConstrainedBox(
          constraints: const BoxConstraints(maxWidth: 420),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              if (failed)
                Container(
                  width: 80,
                  height: 80,
                  decoration: BoxDecoration(color: LumoraColors.coral.withValues(alpha: 0.1), shape: BoxShape.circle),
                  child: const Icon(Icons.cancel, color: LumoraColors.coral, size: 44),
                )
              else
                const FoxMascot(size: 110, glow: true),
              const SizedBox(height: 20),
              Text(title, textAlign: TextAlign.center, style: const TextStyle(fontSize: 18, fontWeight: FontWeight.w800, color: LumoraColors.ink)),
              const SizedBox(height: 4),
              Text(detail, textAlign: TextAlign.center, style: const TextStyle(color: LumoraColors.slatey, height: 1.4)),
              if (actions.isNotEmpty) const SizedBox(height: 24),
              ...actions,
            ],
          ),
        ),
      ),
    );
  }
}
