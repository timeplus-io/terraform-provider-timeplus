FROM alpine:3.19.0

RUN \
    # fetch terraform command line tool binary \
    cd /tmp && \
    wget -O terraform_1.6.5_linux_amd64.zip https://releases.hashicorp.com/terraform/1.6.5/terraform_1.6.5_linux_amd64.zip && \
    unzip terraform_1.6.5_linux_amd64.zip && \
    chmod +x terraform && \
    mv terraform /usr/bin/terraform && \
    #  fetch timeplus provider binary \
    wget -O terraform-provider-timeplus_0.1.4_linux_amd64.zip https://github.com/timeplus-io/terraform-provider-timeplus/releases/download/v0.1.4/terraform-provider-timeplus_0.1.4_linux_amd64.zip && \
    unzip terraform-provider-timeplus_0.1.4_linux_amd64.zip && \
    mkdir -p /terraform/providers && \
    mv terraform-provider-timeplus_v0.1.4 /terraform/providers/terraform-provider-timeplus_v0.1.4 && \
    rm -f * && \
    mkdir /app

WORKDIR /app
VOLUME /app

COPY terraformrc /root/.terraformrc
ENTRYPOINT ["/usr/bin/terraform"]
