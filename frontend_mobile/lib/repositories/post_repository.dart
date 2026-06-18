import 'package:dio/dio.dart';
import '../api/dio_client.dart';
import '../models/post.dart';

class PostRepository {
  final ApiClient _apiClient = ApiClient();

  Future<List<Post>> getFeed({int? limit, int? offset}) async {
    final response = await _apiClient.dio.get('/feed', queryParameters: {
      'limit': limit,
      'offset': offset,
    });
    return (response.data as List).map((json) => Post.fromJson(json)).toList();
  }

  Future<void> createPost({required String content}) async {
    await _apiClient.dio.post('/posts', data: {'content': content});
  }

  Future<void> toggleLike(String postId, bool currentlyLiked) async {
    if (currentlyLiked) {
      await _apiClient.dio.delete('/posts/$postId/like');
    } else {
      await _apiClient.dio.post('/posts/$postId/like');
    }
  }
}