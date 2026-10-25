import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:flutter_application_1/auth/user_model.dart';
import 'package:flutter_application_1/api_client.dart';

class AuthProvider extends ChangeNotifier {
  User? _user;
  bool _isLoading = false;

  User? get user => _user;
  bool get isLoading => _isLoading;
  bool get isLoggedIn => _user != null;

  Future handleLogin(String email, String password) async {
    // final email = _emailController.text;
    // final password = _passwordController.text;

    debugPrint('Email: $email');
    debugPrint('Password: $password');

    try {
        final response = await ApiClient.instance.post(
        '/auth/login',
        data: {
          'email': email,
          'password': password,
        },
      );

      if (response.statusCode == 200) {
        debugPrint('Login successful!');

        // if (mounted) {
        //   Navigator.pop(context);
        // }
      }
      } catch (e) {
        debugPrint('Login error: $e');
      }
  }

  Future<User> fetchUser() async {
    _isLoading = true;
    notifyListeners();

    try {
      final response = await ApiClient.instance.get(
        '/users/me',
      );
      if (response.statusCode == 200) {
        final json = jsonDecode(response) as Map<String, dynamic>;
        return User.fromJson(json);
      }
    } catch (e) {
      _user = null;
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  void clearUser() {
    _user = null;
    notifyListeners();
  }
}