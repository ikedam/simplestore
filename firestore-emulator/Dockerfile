FROM eclipse-temurin:11.0.20.1_1-jre

WORKDIR /workspace
ENV CLOUDSDK_CORE_PROJECT=testproject
ENV PATH=/opt/node/bin:$PATH

RUN curl -sLo /tmp/nodejs.tgz "https://nodejs.org/download/release/v16.20.2/node-v16.20.2-linux-$(uname -m|sed -e 's/86_//').tar.gz" \
  && mkdir -p /opt/node \
  && tar xzf /tmp/nodejs.tgz --strip-components=1 -C /opt/node \
  && rm -rf /tmp/nodejs.tgz \
  && npm install -g firebase-tools@12.5.4 \
  && firebase setup:emulators:firestore \
  && firebase setup:emulators:ui

COPY entrypoint.sh firebase.json /workspace

VOLUME /workspace/data
EXPOSE 8080
EXPOSE 4000

CMD ["sh", "/workspace/entrypoint.sh"]
