import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:frontend_mobile/repositories/post_repository.dart';
import 'package:frontend_mobile/screens/feed_screen.dart';
import 'bloc/feed/feed_bloc.dart';
import 'repositories/auth_repository.dart';
import 'bloc/auth/auth_bloc.dart';
import 'screens/auth_screen.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MultiBlocProvider(
      providers: [
        BlocProvider(create: (context) => AuthBloc(AuthRepository())),
        BlocProvider(create: (context) => FeedBloc(PostRepository())),
      ],
      child: MaterialApp(
        title: 'SocialNet',
        theme: ThemeData(primarySwatch: Colors.blue),
        home: AuthScreen(),
        routes: {
          '/feed': (context) => FeedScreen(),
        },
      ),
    );
  }
}