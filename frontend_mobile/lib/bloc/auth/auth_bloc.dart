import 'package:flutter_bloc/flutter_bloc.dart';
import 'auth_event.dart';
import 'auth_state.dart';
import '../../repositories/auth_repository.dart';

class AuthBloc extends Bloc<AuthEvent, AuthState> {
  final AuthRepository _authRepository;

  AuthBloc(this._authRepository) : super(AuthInitial()) {
    on<LoginRequested>(_onLoginRequested);
    on<RegisterRequested>(_onRegisterRequested);
    on<LogoutRequested>(_onLogoutRequested);
  }

  Future<void> _onLoginRequested(
    LoginRequested event,
    Emitter<AuthState> emit,
  ) async {
    // 1. Проверка текущего состояния для предотвращения дублирования запросов
    if (state is AuthLoading) return;

    emit(AuthLoading());
    try {
      await _authRepository.login(event.username, event.password);
      emit(AuthSuccess());
    } on AuthException catch (e) {
      // 2. Обработка специфичных ошибок репозитория вместо общего catch
      emit(AuthFailure(e.message));
    } catch (e) {
      emit(AuthFailure('Неизвестная ошибка при входе'));
    }
  }

  Future<void> _onRegisterRequested(
    RegisterRequested event,
    Emitter<AuthState> emit,
  ) async {
    if (state is AuthLoading) return;

    emit(AuthLoading());
    try {
      await _authRepository.register(
        event.username,
        event.email,
        event.password,
      );
      emit(AuthSuccess());
    } on AuthException catch (e) {
      emit(AuthFailure(e.message));
    } catch (e) {
      emit(AuthFailure('Неизвестная ошибка при регистрации'));
    }
  }

  Future<void> _onLogoutRequested(
    LogoutRequested event,
    Emitter<AuthState> emit,
  ) async {
    try {
      await _authRepository.logout();
    } catch (_) {
      // 3. Логаут должен сбрасывать состояние даже если запрос к API упал
      // (токен мог быть уже невалиден на сервере)
    } finally {
      emit(AuthInitial());
    }
  }
}

class AuthException implements Exception {
  final String message;
  const AuthException(this.message);
}