cmd/mailman
============

Send mass mails worker.

Make Docker Image
------------------

    sh ./make.sh

Scale
------

Using [docker-compose](https://docs.docker.com/compose/) to scale works.

Setting all values in `./docker-compose.yml`

Scale mailman workers

    docker-compose scale mailman=[nums]

View logs

    docker-compose logs
