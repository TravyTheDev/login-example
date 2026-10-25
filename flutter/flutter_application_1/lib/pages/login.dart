import 'package:flutter/material.dart';
import 'package:flutter_application_1/api_client.dart';
import 'package:flutter_application_1/auth/auth_provider.dart';
import 'package:provider/provider.dart';

class LoginPage extends StatefulWidget {
  const LoginPage({super.key});

  @override
  State<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {
  final GlobalKey<FormState> _formKey = GlobalKey<FormState>();

  final TextEditingController _emailController = TextEditingController();
  final TextEditingController _passwordController = TextEditingController();

  @override
  void dispose() {
    _emailController.dispose();
    _passwordController.dispose();
    super.dispose();
  }

  Future _handleLogin() async {
    final email = _emailController.text;
    final password = _passwordController.text;

    debugPrint('Email: $email');
    debugPrint('Password: $password');

    try {
        await context.read<AuthProvider>().handleLogin(email, password);
        if (mounted) {
          Navigator.pop(context);
        }
      } catch (e) {
        debugPrint('Login error: $e');
      }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Second Route')),
      body: Center(
        child: Column(
          children: [
            Form(
              key: _formKey,
              child: Column(
                crossAxisAlignment: .start,
                children: <Widget>[
                  TextFormField(
                    controller: _emailController,
                    decoration: const InputDecoration(
                      border: OutlineInputBorder(),
                      hintText: 'Enter your email',
                    ),
                  ),
                  TextFormField(
                    controller: _passwordController,
                    obscureText: true,
                    decoration: const InputDecoration(
                      border: OutlineInputBorder(),
                      hintText: 'Enter your password',
                    ),
                  ),
                  Padding(
                    padding: const .symmetric(vertical: 16.0),
                    child: ElevatedButton(
                      onPressed: () {
                        _handleLogin();
                      },
                      child: const Text('Submit'),
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}
