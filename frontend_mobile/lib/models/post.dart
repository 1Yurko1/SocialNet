class Post {
  final String id;
  final String authorId;
  final String? authorName;
  final String content;
  final String? mediaUrl;
  final DateTime createdAt;
  final int likesCount;
  final bool isLiked;

  Post({
    required this.id,
    required this.authorId,
    this.authorName,
    required this.content,
    this.mediaUrl,
    required this.createdAt,
    this.likesCount = 0,
    this.isLiked = false,
  });

  factory Post.fromJson(Map<String, dynamic> json) {
    return Post(
      id: json['id'],
      authorId: json['author_id'],
      authorName: json['author_name'],
      content: json['content'],
      mediaUrl: json['media_url'],
      createdAt: DateTime.parse(json['created_at']),
      likesCount: json['likes_count'] ?? 0,
      isLiked: json['is_liked'] ?? false,
    );
  }
}