FROM rabbitmq:3.13-management

# скачиваем плагин
RUN apt-get update && apt-get install -y curl

RUN curl -L -o /plugins/rabbitmq_delayed_message_exchange.ez \
  https://github.com/rabbitmq/rabbitmq-delayed-message-exchange/releases/download/v3.13.0/rabbitmq_delayed_message_exchange-3.13.0.ez

# включаем плагин
RUN rabbitmq-plugins enable --offline rabbitmq_delayed_message_exchange