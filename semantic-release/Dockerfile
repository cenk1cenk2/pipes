FROM node:16-alpine

RUN yarn global add semantic-release

COPY dist/semantic-release /pipe

RUN chmod +x /pipe

ENTRYPOINT [ "/pipe" ]

