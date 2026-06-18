import 'package:dio/dio.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import '../api/dio_client.dart';
import '../bloc/auth/auth_bloc.dart';

class AuthRepository {
  final ApiClient _apiClient = ApiClient();
  final _storage = const FlutterSecureStorage();

  // Регистрация
  Future<void> register(String username, String email, String password) async {
    await _apiClient.dio.post('auth/register', data: {
      'username': username,
      'email': email,
      'password': password,
    });
  }

  // Вход
  Future<String> login(String username, String password) async {
    final response = await _apiClient.dio.post('/auth/login', data: {
      'username': username,
      'password': password,
    });

    if (response.statusCode == 401)
      throw const AuthException('Неверный логин или пароль');

    String token = response.data['token'];
    await _storage.write(key: 'token', value: token); // Сохраняем JWT
    return token;
  }

  // Проверка: залогинен ли пользователь?
  Future<bool> isAuthenticated() async {
    String? token = await _storage.read(key: 'token');
    return token != null;
  }

  // Выход
  Future<void> logout() async {
    await _storage.delete(key: 'token');
  }
}