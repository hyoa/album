export default {
  fr: {
    auth: {
      connection: 'Connexion',
      register: 'Inscription',
      unavailableWebsite: "Le site est actuellement indisponible ! Mais je travaille d'arrache pied à le remettre en ligne :)",
      alert: {
        info: {
          title: 'Oups'
        },
        error: {
          title: 'Oups'
        },
        disconnect: {
          message: 'Par mesure de sécurité, vous avez été déconnecté. Vous pouvez vous reconnecter avec le formulaire ci-dessous.'
        }
      },
      registerPage: {
        title: "S'inscrire",
        alert: {
          success: {
            title: 'Bienvenue',
            message: 'Votre compte sera bientôt validé et vous pourrez ensuite vous connecter !'
          }
        },
        form: {
          email: 'Email',
          name: 'Nom/Prénom',
          password: 'Mot de passe',
          passwordCheck: 'Confirmation du mot de passe',
          submit: "S'inscrire"
        }
      },
      askResetPasswordPage: {
        alert: {
          success: {
            title: 'Succès',
            message: 'Un email avec la procédure pour réinitialiser votre mot de passe a été envoyé à l\'email indiqué.'
          }
        },
        title: 'Demander un nouveau mot de passe',
        form: {
          email: 'Email',
          submit: 'Demander un nouveau mot de passe'
        }
      },
      loginPage: {
        title: 'Se connecter',
        createAccount: 'Créer un compte maintenant !',
        form: {
          email: 'Email',
          password: 'Mot de passe',
          submit: 'Se connecter'
        },
        forgottenPassword: 'Mot de passe oublié ?'
      }
    },
    userResetPassword: {
      title: 'Réinitialiser mon mot de passe',
      form: {
        password: 'Mot de passe',
        checkPassword: 'Confirmation du mot de passe',
        submit: 'Réinitialiser mon mot de passe'
      }
    },
    home: {
      notificationSection: {
        title: 'Notifications',
        message: 'Vous pouvez désormais recevoir des notifications pour savoir quand de nouveaux albums sont disponibles !',
        advice: "Il est préférable d'installer l'application pour bénéficier au mieux des notifications.",
        help: 'Aide',
        form: {
          accept: 'Accepter de recevoir des notifications',
          refuse: 'Refuser'
        }
      },
      searchFormSection: {
        title: 'Rechercher des albums',
        form: {
          submit: 'Rechercher'
        }
      },
      searchResultSection: {
        title: 'Résultat de la recherche'
      },
      lastAlbumSection: {
        title: 'Les derniers albums',
        loadMore: 'Voir plus d\'album'
      }
    },
    album: {
      createAtBy: 'Créé le {date} par {author}',
      readMore: {
        show: 'Lire la suite',
        hide: 'Cacher'
      }
    },
    error: {
      emailAlreadyUsed: 'Cet adresse email est déjà utilisé !',
      albumAlreadyExist: 'Cet album exist déjà.',
      titleCannotBeEmpty: 'Le titre ne peut pas être vide.',
      default: 'Quelque chose ne s\'est pas déroulé comme il faut. Mais quoi ?!'
    },
    defaultNav: {
      admin: 'Administration',
      disconnect: 'Se déconnecter'
    },
    grid: {
      media: {
        description: 'Photo prise par {author}'
      },
      loading: 'Chargement en cours...'
    },
    admin: {
      adminNav: {
        title: 'Administration',
        leave: 'Quitter',
        album: {
          title: 'Album',
          add: 'Créer',
          see: 'Voir'
        },
        medias: {
          title: 'Médias',
          add: 'Ajouter',
          see: 'Voir'
        },
        users: {
          title: 'Utilisateurs',
          see: 'Voir'
        }
      },
      home: {
        albumCard: {
          title: 'Albums',
          publicCount: 'albums publics',
          privateCount: 'albums privés',
          add: 'Créer',
          see: 'Voir'
        },
        mediaCard: {
          title: 'Médias',
          photos: 'photos',
          videos: 'vidéos',
          add: 'Ajouter',
          see: 'Voir'
        },
        userCard: {
          title: 'Utilisateurs',
          count: 'utilisateurs',
          waitingValidation: 'en attente de validation',
          see: 'Voir'
        }
      },
      albumList: {
        title: 'Liste des albums',
        form: {
          filter: 'Filtrer'
        },
        item: {
          createBy: 'Créer par {author}',
          mediasCount: '{number} médias',
          noMedias: 'Aucun médias'
        }
      },
      albumEdit: {
        form: {
          title: 'Titre',
          description: 'Description',
          private: 'Privé',
          submit: 'Enregister'
        },
        library: 'Bibliothèque',
        mediaSelected: {
          count: '{count} média(s) sélectionné(s)',
          remove: 'Retirer de l\'album'
        },
        notify: {
          editSuccess: 'La modification a été enregistré',
          mediaRemoveSuccess: 'Les médias ont correctement été retirés',
          mediaAddSuccess: 'Les médias ont été ajoutés',
          albumDoesNotExist: 'Cet album n\'existe pas'
        },
        sidebar: {
          title: 'Dossiers',
          close: 'Fermer',
          searchPlaceholder: 'Rechercher un dossier',
          addButton: 'Ajouter',
          selectMedia: 'Sélectionner un ou plusieurs médias'
        }
      },
      mediaFolderList: {
        title: 'Liste des dossiers'
      },
      mediaFolder: {
        form: {
          folderName: 'nom du dossier',
          submit: 'Mettre à jour'
        },
        mediaSelected: {
          count: '{count} média(s) sélectioné(s)',
          form: {
            folder: 'Dossier',
            folderPlaceholder: 'Nouveau dossier',
            submit: 'Changer de dossier'
          }
        },
        notify: {
          submitSuccess: 'Enregistrement effectué',
          moveSuccess: 'Les médias sélectionnés ont correctement été déplacés'
        }
      },
      mediaAdd: {
        title: 'Ajouter des médias',
        uploadRunning: {
          title: 'Téléchargement en cours...',
          timeRemaining: 'Temps restant: {time}',
          totalUploaded: 'Total: {count}',
          successUploaded: 'Réussi: {count}',
          failUploaded: 'Echec: {count}',
          timeRemaingingCalculation: 'Calcul du temps restant en cours',
          almostDone: 'C\'est bientôt terminé !'
        },
        form: {
          folder: 'Dossier',
          media: 'Médias',
          dragAndDrop: 'Clique ou dépose tes fichiers ici',
          album: 'Album',
          linkToAlbum: 'Lié à un album'
        },
        notify: {
          tooManyMedia: 'Il n\'est pas possible d\'envoyer plus de 20 fichier à la fois',
          uploadSuccess: 'Les fichiers ont été transférés sur le serveur'
        }
      },
      userList: {
        title: 'Liste des utilisateurs',
        sectionInvitation: {
          title: 'Envoyer une invitation',
          submit: 'Inviter'
        },
        sectionList: {
          table: {
            name: 'Nom',
            role: 'Role',
            action: 'Actions',
            activate: 'Activer'
          }
        },
        notify: {
          accountActivate: 'Le compte a été activé',
          invitationSend: 'Invitations envoyées'
        }
      }
    },
    components: {
      autocomplete: {
        isNew: 'nouveau'
      }
    }
  }
}
