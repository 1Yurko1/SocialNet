import 'package:dio/dio.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class ApiClient {
  static const String _baseUrl = 'http://192.168.1.71:8080/api/v1';

  final Dio dio;
  final FlutterSecureStorage _storage;

  ApiClient({
    Dio? dio,
    FlutterSecureStorage? storage,
  })  : dio = dio ?? Dio(BaseOptions(
          baseUrl: _baseUrl,
          connectTimeout: const Duration(seconds: 5),
          receiveTimeout: const Duration(seconds: 10), // 4. Увеличен таймаут чтения
        )),
        _storage = storage ?? const FlutterSecureStorage() {
    _setupInterceptors();
  }

  void _setupInterceptors() {
    dio.interceptors.add(InterceptorsWrapper(
      onRequest: (options, handler) async {
        final token = await _storage.read(key: 'token');
        if (token != null && token.isNotEmpty) {
          options.headers['Authorization'] = 'Bearer $token';
        }
        handler.next(options);
      },
      onError: (DioException e, handler) async {
        if (e.response?.statusCode == 401) {
          // TODO: Добавить логику обновления токена (refresh token)
          // Если рефреш не удался — выполнить выход из аккаунта

          // Пример: предотвращение бесконечного цикла при ошибке рефреша
          // if (!_isRefreshing) { ... }
        }
        handler.next(e);
      },
    ));
  }
}