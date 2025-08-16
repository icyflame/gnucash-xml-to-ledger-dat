FROM perl:5.42

ENV PERL_MM_USE_DEFAULT 1
ENV PERL5LIB /home/lib/perl5
ENV PERL_MB_OPT --install_base /home
ENV PERL_MM_OPT INSTALL_BASE /home

RUN cpan -f -i App:cpanminus
RUN cpanm --notest XML::Simple

COPY copy.sh convert.pl /src/

ENTRYPOINT ["bash", "/src/copy.sh"]