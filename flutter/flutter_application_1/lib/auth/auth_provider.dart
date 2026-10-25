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
    debugPrint('Email: $email');
    debugPrint('Password: $password');
    _isLoading = true;
    try {
      final response = await ApiClient.instance.post(
        '/auth/login',
        data: {'email': email, 'password': password},
      );
      if (response.statusCode == 200) {
        final json = response.data as Map<String, dynamic>;
        final user = User.fromJson(json);

        _user = user;
        return user;
      }
    } catch (e) {
      debugPrint('Login error: $e');
    } finally {
      _isLoading = false;
    }
  }

  Future<User> fetchUser() async {
    _isLoading = true;
    notifyListeners();

    try {
      final response = await ApiClient.instance.get('/users/me');

      if (response.statusCode == 200) {
        final json = response.data as Map<String, dynamic>;
        final user = User.fromJson(json);

        _user = user;
        return user;
      }

      throw Exception('Failed to fetch user');
    } catch (e) {
      _user = null;
      rethrow;
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
