# $NetBSD: Makefile,v 1.19 2018/08/22 09:44:07 wiz Exp $

DISTNAME=	IO-CaptureOutput-1.1104
PKGNAME=	p5-${DISTNAME:S/1104/11.04/}
PKGREVISION=	4
CATEGORIES=	devel perl5
MASTER_SITES=	${MASTER_SITE_PERL_CPAN:=IO/}

MAINTAINER=	pkgsrc-users@NetBSD.org
HOMEPAGE=	http://search.cpan.org/dist/IO-CaptureOutput/
COMMENT=	Capture STDOUT and STDERR from Perl code, subprocesses, or XS
LICENSE=	${PERL5_LICENSE}

PERL5_PACKLIST=		auto/IO/CaptureOutput/.packlist

#DEPENDS+=	{perl>=5.8.9,p5-PathTools>=3.27}:../../devel/p5-PathTools

.include "../../lang/perl5/module.mk"
.include "../../mk/bsd.pkg.mk"
