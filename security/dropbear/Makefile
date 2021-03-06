# $NetBSD: Makefile,v 1.35 2018/07/04 13:40:33 jperkin Exp $

DISTNAME=	dropbear-2017.75
CATEGORIES=	security
MASTER_SITES=	http://matt.ucc.asn.au/dropbear/releases/
EXTRACT_SUFX=	.tar.bz2

MAINTAINER=	snj@NetBSD.org
HOMEPAGE=	http://matt.ucc.asn.au/dropbear/dropbear.html
COMMENT=	Small SSH2 server and client, aimed at embedded market
LICENSE=	modified-bsd

GNU_CONFIGURE=		yes
CONFIGURE_ARGS+=	--sysconfdir=${PKG_SYSCONFDIR:Q} --disable-bundled-libtom
USE_TOOLS+=		gmake

PKG_OPTIONS_VAR=	PKG_OPTIONS.dropbear
PKG_SUPPORTED_OPTIONS=	pam

.include "../../mk/bsd.prefs.mk"
.include "../../mk/bsd.options.mk"

.if !empty(PKG_OPTIONS:Mpam)
.  include "../../mk/pam.buildlink3.mk"
CONFIGURE_ARGS+=	--enable-pam
SUBST_CLASSES+=		pam
SUBST_MESSAGE.pam=	Enabling PAM in options.h
SUBST_STAGE.pam=	pre-configure
SUBST_FILES.pam=	options.h
SUBST_SED.pam=		-e "s/ENABLE_SVR_PASSWORD_AUTH/ENABLE_SVR_PAM_AUTH/"
.endif

OWN_DIRS+=		${PKG_SYSCONFDIR}/dropbear

SUBST_CLASSES+=		config
SUBST_MESSAGE.config=	Fixing path to config directory.
SUBST_STAGE.config=	post-build
SUBST_FILES.config=	dropbear.8 dropbearkey.1
SUBST_SED.config=	-e "s,/etc/dropbear/,"${PKG_SYSCONFDIR:Q}"/dropbear/,g"

# used by dbscp
CPPFLAGS+=		-DDROPBEAR_PATH_SSH_PROGRAM="\"${PREFIX}/bin/dbclient\""

.include "../../x11/xauth/builtin.mk"

CPPFLAGS+=		-DXAUTH_COMMAND="\"${XAUTHBASE}/bin/xauth\""

CFLAGS.NetBSD+=		-DHAVE_NETINET_IN_SYSTM_H
LDFLAGS.SunOS+=		-lsocket -lnsl

INSTALLATION_DIRS=	share/doc/dropbear ${PKGMANDIR}/man1 ${PKGMANDIR}/man8

BUILD_TARGET=		all scp

post-install:
	${INSTALL_PROGRAM} ${WRKSRC}/scp ${DESTDIR}/${PREFIX}/bin/dbscp

.include "../../devel/zlib/buildlink3.mk"
.include "../../math/ltm/buildlink3.mk"
.include "../../security/libtomcrypt/buildlink3.mk"
.include "../../mk/bsd.pkg.mk"
