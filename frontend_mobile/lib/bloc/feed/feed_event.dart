abstract class FeedEvent {}

class FetchPostsRequested extends FeedEvent {}
class CreatePostRequested extends FeedEvent {
  final String content;
  CreatePostRequested(this.content);
}
class ToggleLikeRequested extends FeedEvent {
  final String postId;
  final bool isLiked;
  ToggleLikeRequested(this.postId, this.isLiked);
}