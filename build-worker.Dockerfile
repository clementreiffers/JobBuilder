FROM ubuntu AS worker-builder

RUN apt-get update && apt-get install -y clang libc++-dev nodejs npm

COPY ./ ./

RUN npx workerd compile config.capnp > serv.out

FROM ubuntu AS worker

RUN apt-get update && apt-get install -y libc++-dev

COPY --from=worker-builder serv.out .

CMD ["./serv.out"]
