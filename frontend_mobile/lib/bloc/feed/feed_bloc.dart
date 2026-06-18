import 'package:flutter_bloc/flutter_bloc.dart';
import '../../models/post.dart';
import 'feed_event.dart';
import 'feed_state.dart';
import '../../repositories/post_repository.dart';

class FeedBloc extends Bloc<FeedEvent, FeedState> {
  final PostRepository postRepo;

  FeedBloc(this.postRepo) : super(const FeedState()) {
    on<FetchPostsRequested>((event, emit) async {
      emit(state.copyWith(isLoading: true));
      try {
        final posts = await postRepo.getFeed(limit: 20, offset: 0);
        emit(state.copyWith(posts: posts, isLoading: false));
      } catch (e) {
        emit(state.copyWith(isLoading: false));
      }
    });

    on<CreatePostRequested>((event, emit) async {
      try {
        await postRepo.createPost(content: event.content);
        // После создания просто обновляем ленту
        add(FetchPostsRequested());
      } catch (e) {
        // Обработка ошибки
      }
    });

    on<ToggleLikeRequested>((event, emit) async {
      try {
        await postRepo.toggleLike(event.postId, event.isLiked);
        // Оптимистичное обновление списка
        final updatedPosts = state.posts.map((p) {
          if (p.id == event.postId) {
            return Post(
              id: p.id,
              authorId: p.authorId,
              authorName: p.authorName,
              content: p.content,
              createdAt: p.createdAt,
              isLiked: !event.isLiked,
              likesCount: p.likesCount + (event.isLiked ? -1 : 1),
            );
          }
          return p;
        }).toList();
        emit(state.copyWith(posts: updatedPosts));
      } catch (e) {
        // Ошибка
      }
    });
  }
}
