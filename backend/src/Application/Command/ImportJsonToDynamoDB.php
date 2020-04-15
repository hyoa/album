<?php

namespace Album\Application\Command;

use Album\Domain\Album\AlbumEntity;
use Album\Domain\Album\AlbumMediaEntity;
use Album\Domain\Album\AlbumRepositoryInterface;
use Album\Domain\Media\MediaEntity;
use Album\Domain\Media\MediaRepositoryInterface;
use Album\Domain\User\UserEntity;
use Album\Domain\User\UserRepositoryInterface;
use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputArgument;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;
use Symfony\Component\Console\Style\SymfonyStyle;

class ImportJsonToDynamoDB extends Command
{
    protected static $defaultName = 'app:dynamodb:import';

    protected AlbumRepositoryInterface $albumRepository;

    protected MediaRepositoryInterface $mediaRepository;

    protected UserRepositoryInterface $userRepository;

    public function __construct(
        AlbumRepositoryInterface $albumRepository,
        MediaRepositoryInterface $mediaRepository,
        UserRepositoryInterface $userRepository,
        string $name = null
    ) {
        parent::__construct($name);

        $this->albumRepository = $albumRepository;
        $this->mediaRepository = $mediaRepository;
        $this->userRepository = $userRepository;
    }

    public function configure()
    {
        $this->addArgument('path', InputArgument::REQUIRED);
        $this->addArgument('table', InputArgument::REQUIRED);
    }

    public function execute(InputInterface $input, OutputInterface $output)
    {
        $io = new SymfonyStyle($input, $output);

        $io->title('Import JSON to DynamoDB');

        $table = $input->getArgument('table');

        $file = fopen($input->getArgument('path'), 'r');

        $data = [];
        if ($file) {
            while (($buffer = fgets($file)) !== false) {
                $data[] = json_decode($buffer, true);
            }
            if (!feof($file)) {
                echo "Erreur: fgets() a échoué\n";
            }
            fclose($file);
        }

        switch ($table) {
            case 'user':
                $this->importUsers($data);
                break;
            case 'album':
                $this->importAlbums($data);
                break;
            case 'media':
                $this->importMedias($data);
                break;
            default:
                throw new \Exception();
        }

        return 1;
    }

    protected function importUsers(array $data)
    {
        foreach ($data as $user) {
            $userEntity = new UserEntity();
            $userEntity->name = $user['name'];
            $userEntity->email = $user['email'];
            $userEntity->password = $user['password'];
            $userEntity->role = $user['role']['$numberInt'];
            $userEntity->registrationDate = (new \DateTimeImmutable())->setTimestamp($user['registrationDate']['$numberInt']);

            $this->userRepository->insert($userEntity);
        }
    }

    protected function importMedias(array $data)
    {
        foreach ($data as $media) {
            $mediaEntity = new MediaEntity();
            $mediaEntity->type = $media['type']['$numberInt'];
            $mediaEntity->key = $media['key'];
            $mediaEntity->author = $media['author'];
            $mediaEntity->uploadDate = (new \DateTimeImmutable())->setTimestamp($media['uploadDate']['$numberInt']);
            $mediaEntity->folder = $media['folder'];

            $this->mediaRepository->insert($mediaEntity);
        }
    }

    protected function importAlbums(array $data)
    {
        foreach ($data as $album) {
            $albumEntity = new AlbumEntity();
            $albumEntity->title = $album['title'];
            $albumEntity->description = !empty($album['description']) ? $album['description'] : null;
            $albumEntity->author = $album['author'];
            $albumEntity->private = $album['private'];
            $albumEntity->creationDate = (new \DateTimeImmutable())->setTimestamp($album['creationDate']['$numberInt']);
            $albumEntity->slug = $album['slug'];


            if (isset($album['medias'])) {
                dump('d');
                foreach ($album['medias'] as $media) {
                    $mediaEntity = new AlbumMediaEntity();
                    $mediaEntity->author = $media['author'];
                    $mediaEntity->key = $media['key'];
                    $mediaEntity->type = $media['type']['$numberInt'];
                    $mediaEntity->isFavorite = $media['favorite'] ?? false;

                    $albumEntity->addMedia($mediaEntity);
                }
            }


            try {
                $this->albumRepository->insert($albumEntity);
            } catch (\Exception $exception) {
                dump($album);
            }
        }
    }
}