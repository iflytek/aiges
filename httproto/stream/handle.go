package stream

type Handler func(s *WsSession)

type HandlerChain []Handler
