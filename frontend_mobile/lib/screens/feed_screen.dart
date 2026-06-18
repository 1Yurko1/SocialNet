import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../bloc/feed/feed_bloc.dart';
import '../bloc/feed/feed_event.dart';
import '../bloc/feed/feed_state.dart';
import '../models/post.dart';

class FeedScreen extends StatefulWidget {
  @override
  _FeedScreenState createState() => _FeedScreenState();
}

class _FeedScreenState extends State<FeedScreen> {
  final _postController = TextEditingController();

  @override
  void initState() {
    super.initState();
    // Загружаем посты при открытии экрана
    context.read<FeedBloc>().add(FetchPostsRequested());
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.grey[100],
      appBar: AppBar(
        title: Text("SocialNet", style: TextStyle(fontWeight: FontWeight.bold, color: Colors.blue[600])),
        backgroundColor: Colors.white,
        elevation: 0,
        actions: [
          IconButton(
            icon: Icon(Icons.logout, color: Colors.grey),
            onPressed: () => Navigator.of(context).pushReplacementNamed('/auth'),
          )
        ],
      ),
      body: Column(
        children: [
          // Поле создания поста
          Padding(
            padding: const EdgeInsets.all(16.0),
            child: Container(
              padding: EdgeInsets.symmetric(horizontal: 12, vertical: 8),
              decoration: BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.circular(20),
                boxShadow: [BoxShadow(color: Colors.black12, blurRadius: 4)],
              ),
              child: Row(
                children: [
                  CircleAvatar(backgroundColor: Colors.blue[300]),
                  SizedBox(width: 12),
                  Expanded(
                    child: TextField(
                      controller: _postController,
                      decoration: InputDecoration(
                        hintText: "What's happening?",
                        border: InputBorder.none,
                      ),
                    ),
                  ),
                  IconButton(
                    icon: Icon(Icons.send, color: Colors.blue[600]),
                    onPressed: () {
                      if (_postController.text.isNotEmpty) {
                        context.read<FeedBloc>().add(CreatePostRequested(_postController.text));
                        _postController.clear();
                      }
                    },
                  ),
                ],
              ),
            ),
          ),

          // Список постов
          Expanded(
            child: BlocBuilder<FeedBloc, FeedState>(
              builder: (context, state) {
                if (state.isLoading && state.posts.isEmpty) {
                  return Center(child: CircularProgressIndicator());
                }
                if (state.posts.isEmpty) {
                  return Center(child: Text("No posts yet!"));
                }

                return ListView.builder(
                  padding: EdgeInsets.symmetric(horizontal: 16),
                  itemCount: state.posts.length,
                  itemBuilder: (context, index) {
                    final post = state.posts[index];
                    return PostCard(post: post);
                  },
                );
              },
            ),
          ),
        ],
      ),
    );
  }
}

// Отдельный виджет карточки поста для чистоты кода
class PostCard extends StatelessWidget {
  final Post post;
  const PostCard({required this.post});

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: EdgeInsets.only(bottom: 16),
      padding: EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(24),
        border: Border.all(color: Colors.cyan[200] ?? Colors.grey[200]!),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              CircleAvatar(backgroundColor: Colors.grey[300]),
              SizedBox(width: 12),
              Text(post.authorName ?? "User", style: TextStyle(fontWeight: FontWeight.bold)),
              Spacer(),
              Text("Just now", style: TextStyle(fontSize: 12, color: Colors.grey)),
            ],
          ),
          SizedBox(height: 12),
          Text(post.content, style: TextStyle(fontSize: 16, color: Colors.black87)),
          SizedBox(height: 12),
          Row(
            children: [
              GestureDetector(
                onTap: () => context.read<FeedBloc>().add(ToggleLikeRequested(post.id, post.isLiked)),
                child: Row(
                  children: [
                    Icon(
                      Icons.favorite,
                      color: post.isLiked ? Colors.red : Colors.grey,
                      size: 20
                    ),
                    SizedBox(width: 4),
                    Text("${post.likesCount}", style: TextStyle(fontSize: 12)),
                  ],
                ),
              ),
              SizedBox(width: 20),
              Row(
                children: [
                  Icon(Icons.chat_bubble_outline, color: Colors.grey, size: 20),
                  SizedBox(width: 4),
                  Text("Comment", style: TextStyle(fontSize: 12, color: Colors.grey)),
                ],
              ),
            ],
          ),
        ],
      ),
    );
  }
}