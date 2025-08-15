package redisctldto

type GetByIDRedis struct {
	ID string `uri:"id" binding:"required"`
}
