package logging

// func NewLogger(logger log.Logger) jet.Middleware {
// 	return func(next jet.Handler) jet.Handler {
// 		return func(ctx context.Context, name string, request any) (response any, err error) {
// 			defer func() {
// 				level := log.LevelInfo
// 				if err != nil {
// 					level = log.LevelError
// 				}
//
// 				log.NewHelper(log.WithContext(ctx, logger)).
// 					Log(level, "name", name, "request", request, "response", response, "error", err)
// 			}()
// 			return next(ctx, name, request)
// 		}
// 	}
// }
