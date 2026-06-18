import 'package:flutter_bloc/flutter_bloc.dart';
import 'auth_event.dart';
import 'auth_state.dart';
import '../../repositories/auth_repository.dart';

class AuthBloc extends Bloc<AuthEvent, AuthState> {
  final AuthRepository authRepository;

  AuthBloc(this.authRepository) : super(AuthInitial()) {
    // Обработка входа
    on<LoginRequested>((event, emit) async {
      emit(AuthLoading());
      try {
        await authRepository.login(event.username, event.password);
        emit(AuthSuccess());
      } catch (e) {
        emit(AuthFailure("Login failed. Check your credentials."));
      }
    });

    // Обработка регистрации
    on<RegisterRequested>((event, emit) async {
      emit(AuthLoading());
      try {
        await authRepository.register(event.username, event.email, event.password);
        emit(AuthSuccess()); // Или отдельное состояние RegisterSuccess
      } catch (e) {
        emit(AuthFailure("Registration failed."));
      }
    });

    on<LogoutRequested>((event, emit) async {
      await authRepository.logout();
      emit(AuthInitial());
    });
  }
}