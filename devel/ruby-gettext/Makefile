# $NetBSD: Makefile,v 1.35 2018/03/13 16:49:03 taca Exp $

DISTNAME=	gettext-3.2.9
CATEGORIES=	devel

MAINTAINER=	pkgsrc-users@NetBSD.org
HOMEPAGE=	http://ruby-gettext.github.com/
COMMENT=	Ruby localization library like as GNU gettext
LICENSE=	ruby-license OR gnu-lgpl-v3

DEPENDS+=	${RUBY_PKGPREFIX}-locale>=2.0.5:../../devel/ruby-locale
DEPENDS+=	${RUBY_PKGPREFIX}-text>=1.3.0:../../textproc/ruby-text

RUBYGEM_OPTIONS+=	--format-executable

.include "../../lang/ruby/gem.mk"
.include "../../mk/bsd.pkg.mk"
