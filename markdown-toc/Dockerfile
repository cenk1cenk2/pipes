FROM node:16-alpine

RUN yarn global add markdown-toc

COPY dist/markdown-toc /pipe

RUN chmod +x /pipe

ENTRYPOINT [ "/pipe" ]
