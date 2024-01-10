# NoobSocial

[![wakatime](https://wakatime.com/badge/user/7bd238cb-c7ea-4e56-abe2-0b6ae36ff252/project/018b7110-76e1-42c4-8efc-b616eb7f6a4b.svg)](https://wakatime.com/badge/user/7bd238cb-c7ea-4e56-abe2-0b6ae36ff252/project/018b7110-76e1-42c4-8efc-b616eb7f6a4b)

The Social Network that aims to be a self hosted network for just you and your friends.

No Algorithms, No Ads, No BS.

Self Hosted, Open Source, and Free. Not only that, it is written using the latest and greatest technologies that are designed to be as pocket friendly as possible and as fast as possible.

So what are you waiting for? Get started today!

## Under Development

I am making this project to have a social media platform for just me and my friends. However, one more reason is to use the Go and HTMX stack to make a fast and efficient social media platform that is also easy to deploy and maintain.

## Setting up for Development

### Prerequisites

- Go 1.2+
- NodeJS 18+
- [Templ](https://github.com/a-h/templ)
- [TailwindCSS](https://tailwindcss.com/)
- [PostgreSQL](https://www.postgresql.org/)

You need two terminals open, one for the backend and one for active tailwind compilation.

For the backend, just open the terminal and run `make run`. This will start the backend server on port 8000 or if you have [air](https://github.com/cosmtrek/air) installed, it will start the server on port 8000 and will automatically reload the server when you make changes.

On the other terminal, use `make tailwind` to start the tailwind compilation. This will watch the `tailwind.css` file and will automatically compile it to `/static/output_prod_styles.css` when you make changes.

## Check List

- [x] Basic Backend
- [x] Authentication
- [ ] Post Creation Frontend
- [ ] Post Management Frontend

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

## Contributing

Feel free to contribute. There is no established contributing guidelines yet, but I will be adding them soon. Until then, just make a pull request and I will review it or open an issue and I will look into it.
