module gglmm-tencentyun-example

go 1.13

replace github.com/weihongguo/gglmm => ../../gglmm

replace github.com/weihongguo/gglmm-redis => ../../gglmm-redis

replace github.com/weihongguo/gglmm-tencentyun => ../

require (
	github.com/weihongguo/gglmm v0.0.0-20200226150144-384f169aa64a
	github.com/weihongguo/gglmm-redis v0.0.0-00010101000000-000000000000
	github.com/weihongguo/gglmm-tencentyun v0.0.0-20200331141323-f4b9922f09e6
)
