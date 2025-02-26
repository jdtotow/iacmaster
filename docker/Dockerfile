FROM bitnami/azure-cli:2.69.0-debian-12-r3
USER root
RUN apt-get update
RUN apt-get install -y git curl unzip  
RUN git clone --depth=1 https://github.com/tfutils/tfenv.git ~/.tfenv
RUN ln -s ~/.tfenv/bin/* /usr/local/bin
ENV PATH=$PATH:/usr/local/bin
RUN mkdir /app 
COPY ./entrypoint.sh /app/
USER 1001