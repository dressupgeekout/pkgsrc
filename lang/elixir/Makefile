# $NetBSD: Makefile,v 1.8 2018/09/04 13:49:25 fhajny Exp $

DISTNAME=		elixir-1.7.3
CATEGORIES=		lang
MASTER_SITES=		${MASTER_SITE_GITHUB:=elixir-lang/}

MAINTAINER=		pkgsrc-users@NetBSD.org
HOMEPAGE=		https://elixir-lang.org/
COMMENT=		Functional, meta-programming aware language built on top of Erlang VM
LICENSE=		apache-2.0

GITHUB_TAG=		v${PKGVERSION_NOREV}

DEPENDS+=		erlang>=18.0:../../lang/erlang

USE_TOOLS+=		gmake

REPLACE_INTERPRETER+=	elixir
REPLACE.elixir.old=	.*elixir
REPLACE.elixir.new=	${PREFIX}/bin/elixir
REPLACE_FILES.elixir=	bin/mix

BUILD_TARGET=		compile

.include "../../mk/bsd.pkg.mk"
