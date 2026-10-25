import 'package:dio/dio.dart';
import 'package:cookie_jar/cookie_jar.dart';
import 'package:dio_cookie_manager/dio_cookie_manager.dart';

class ApiClient {
  static final Dio _dio = Dio(BaseOptions(
    baseUrl: 'http://localhost:8000',
    contentType: 'application/json',
    connectTimeout: const Duration(seconds: 10),
    receiveTimeout: const Duration(seconds: 10),
  ));
    static bool _isInitialized = false;

  static Dio get instance {
    if (!_isInitialized) {
      _dio.interceptors.add(CookieManager(CookieJar()));
      
      _dio.interceptors.add(TokenRefreshInterceptor(_dio));
      
      _isInitialized = true;
    }
    return _dio;
  }

}
class TokenRefreshInterceptor extends QueuedInterceptor {
  final Dio dio;

  TokenRefreshInterceptor(this.dio);

  @override
  Future<void> onError(DioException err, ErrorInterceptorHandler handler) async {
    final path = err.requestOptions.path;

    if (err.response?.statusCode == 401 && !path.contains('/auth/renew_token')) {
      try {
        final refreshResponse = await dio.post('/auth/renew_token');

        if (refreshResponse.statusCode == 200 || refreshResponse.statusCode == 204) {
          final opts = err.requestOptions;
          final response = await dio.fetch(opts);
          
          return handler.resolve(response);
        } else {
          throw DioException(
            requestOptions: err.requestOptions,
            error: 'Session expired',
          );
        }
      } catch (refreshError) {
        return handler.reject(err);
      }
    }

    return super.onError(err, handler);
  }
}