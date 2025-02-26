FROM ubuntu:24.10
RUN apt-get update && apt-get install -y git curl unzip 
RUN curl -sL https://aka.ms/InstallAzureCLIDeb | bash
RUN git clone --depth=1 https://github.com/tfutils/tfenv.git ~/.tfenv
RUN ln -s ~/.tfenv/bin/* /usr/local/bin
ENV PATH=$PATH:/usr/local/bin
RUN mkdir /app 
COPY ./entrypoint.sh /app/