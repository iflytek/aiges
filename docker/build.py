import os.path
import pathlib
from typing import Any

from jinja2 import Environment
from plumbum import cli
from plumbum.cmd import rm  # type: ignore

from config import *

log = logging.getLogger()


class Manager(cli.Application):
    """Aiges CI Manager"""

    PROGNAME: str = "build.py"
    VERSION: str = "0.0.1"

    manifest = {}
    ci = None

    def main(self):
        if not self.nested_command:  # will be ``None`` if no sub-command follows
            log.fatal("No subcommand given!")
            print()
            self.help()
            return 1
        elif len(self.nested_command[1]) < 2 and any(
                "generate" in arg for arg in self.nested_command[1]
        ):
            log.error(
                "Subcommand 'generate' missing  required arguments! use 'generate --help'"
            )
            return 1


class ImageTag(object):
    def __init__(self, chip, python="3.9.13", golang="1.17", distro="ubuntu1804"):
        self.chip = chip
        self.python = python
        self.golang = golang
        self.distro = distro

    def __str__(self):
        return "{chip}-{golang}-{python}-{distro}".format(chip=self.chip, golang=self.golang, python=self.python,
                                                          distro=self.distro)


@Manager.subcommand("generate")  # type: ignore
class ManagerGenerate(Manager):
    DESCRIPTION = "Generate Dockerfiles from templates."

    parent: Manager

    vars = {}

    matrix = []
    template_env: Any = Environment(
        extensions=["jinja2.ext.do", "jinja2.ext.loopcontrols"],
        trim_blocks=True,
        lstrip_blocks=True,
    )

    template: Any

    generate_all: Any = cli.Flag(
        ["--all"],
        help="cpu&&gpu version.",
    )

    use_github: Any = cli.Flag(
        ["--use_github"],
        help="If Using Github Actions",
    )

    use_conda: Any = cli.Flag(
        ["--use_conda"],
        help="If Using Miniconda",
        default=True,
    )

    generate_cpu: Any = cli.Flag(
        ["--cpu"],
        help="cpu version",
    )

    generate_gpu: Any = cli.Flag(
        ["--gpu"],
        help="gpu version",
    )

    py_version: Any = cli.SwitchAttr(
        "--python_version",
        str,
        help="specify python version.",
        default="3.9",
    )

    distro: Any = cli.SwitchAttr(
        "--os-name",
        str,
        help="The distro to use.",
        default="ubuntu",
    )

    distro_version: Any = cli.SwitchAttr(
        "--os-version",
        str,
        group="Targeted",
        excludes=["--all"],
        help="The distro version",
        default=None,
    )

    cuda_version: Any = cli.SwitchAttr(
        "--cuda-version",
        str,
        excludes=["--all"],
        group="Targeted",
        help="The cuda version to use. Example: '11.2'",
        default=None,
    )

    action: Any = cli.SwitchAttr(
        "--action",
        str,
        group="Targeted",
        excludes=["--all", ],
        help="Action for build.py 'build' or 'release'",
        default='build',
    )

    git_tag: Any = cli.SwitchAttr(
        "--git_tag",
        str,
        group="Targeted",
        excludes=["--all", ],
        help="git tag ",
        default='v1.2.0',
    )

    def matched(self, key):
        match = self.cuda_version_regex.match(key)
        if match:
            return match

    # extracts arbitrary keys and inserts them into the templating context
    def extract_keys(self, val, arch=None):
        pass

    # For cudnn templates, we need a custom template context
    def output_cudnn_template(self, cudnn_version_name, input_template, output_path):
        pass

    def prepare_context(self):

        # The templating context. This data structure is used to fill the templates.
        self.vars = {
            "registry": self.get_regsitry(),
            "tag": self.generate_matrix_tags(),
        }

    def generate_matrix_tags(self):
        if self.generate_cpu:
            CHIP_LIST = SUPPORTED_CPU_LIST
        elif self.generate_gpu:
            CHIP_LIST = SUPPORTED_CUDA_LIST
        elif self.generate_all:
            CHIP_LIST = SUPPORTED_CPU_LIST + SUPPORTED_CUDA_LIST

        for chip in CHIP_LIST:
            for python in SUPPORTED_PYVERSION_LIST:
                for golang in SUPPORTED_GOLANG_LIST:
                    for distro in SUPPORTED_DISTRO_LIST:
                        self.matrix.append(ImageTag(chip=chip, python=python, golang=golang, distro=distro))


    def generate_release_note(self):
        path = './hack/release/Note.md'
        release_line_format = "| {registry}/{repo}:{tag}{git_tag} | {tag} | {python} | {cuda} | {distro} |"

        base_images_list = [
            release_line_format.format(registry=self.get_regsitry(), tag=str(tag), python=tag.python, cuda=tag.cuda,
                                       distro=tag.distro, repo='cuda-go-python-base', git_tag="")
            for tag in self.matrix]
        log.info(base_images_list)
        aiges_images_list = [
            release_line_format.format(registry=self.get_regsitry(), tag=str(tag), python=tag.python, cuda=tag.cuda,
                                       distro=tag.distro, repo='aiges-gpu', git_tag="-{}".format(self.git_tag))
            for tag in self.matrix]
        log.info(aiges_images_list)

        s = self.release_note.render(vars={
            "base_images": '\n'.join(base_images_list),
            "aiges_images": '\n'.join(aiges_images_list),
        })
        log.info(s)
        with open(path, 'w') as note:
            note.write(s)
            note.close()

    def generate_dockerfile(self):
        if not os.path.exists(TEMP_GEN_DIR):
            os.makedirs(TEMP_GEN_DIR)
        for tag in self.matrix:
            if tag.chip == "cpu":
                if self.use_conda:
                    if self.distro == "ubuntu":
                        self._load_template("./docker/templates/cpu/miniconda/ubuntu/Dockerfile.j2")
                        dockerfile_dir = os.path.join(TEMP_GEN_DIR, tag.distro,tag.chip,"miniconda",self.distro)  
                    elif self.distro == "debian":
                        self._load_template("./docker/templates/cpu/miniconda/debian/Dockerfile.j2")
                        dockerfile_dir = os.path.join(TEMP_GEN_DIR, tag.distro,tag.chip,"miniconda",self.distro)  
                    else:
                        log.error("%s that do not support building Dockerfiles" % self.distro)
                else:
                    if self.distro == "ubuntu":
                        self._load_template("./docker/templates/cpu/ubuntu/Dockerfile.j2")
                        dockerfile_dir = os.path.join(TEMP_GEN_DIR, tag.distro,tag.chip,self.distro)  
                    elif self.distro == "debian":
                        self._load_template("./docker/templates/cpu/debian/Dockerfile.j2")
                        dockerfile_dir = os.path.join(TEMP_GEN_DIR, tag.distro,tag.chip,self.distro)  
                    else:
                        log.error("%s that do not support building Dockerfiles" % self.distro)
            else:
                self._load_template("./docker/templates/aiges-gpu/Dockerfile.j2")
                dockerfile_dir = os.path.join(TEMP_GEN_DIR, tag.distro,"cuda-"+ tag.chip)  
            st = self.render(tag)
            if not os.path.exists(dockerfile_dir):
                os.makedirs(dockerfile_dir)
            with open(os.path.join(dockerfile_dir, Dockerfile), 'w') as dockerfile:
                dockerfile.write(st)
                dockerfile.close()
                log.info("write %s success" % os.path.abspath(os.path.join(dockerfile_dir, Dockerfile)))

    def render(self, tag: ImageTag):
        s = self.template.render(use_github=self.use_github, vars={
                "registry": self.get_regsitry(),
                "tag": str(tag),
                "python_version": self.py_version
            })
        return s

    def get_regsitry(self):
        if self.use_github:
            return ECR_REPO
        return INNER_REPO

    def set_output_path(self, target):
        self.output_path = pathlib.Path(
            f"{self.dist_base_path}/{target.replace('.', '')}"
        )
        if not self.parent.shipit_uuid and self.output_path.exists:
            log.info(f"Removing {self.output_path}")
            rm["-rf", self.output_path]()
        log.debug(f"self.output_path: '{self.output_path}' target: '{target}'")
        log.debug(f"Creating {self.output_path}")
        self.output_path.mkdir(parents=True, exist_ok=False)

    def cheak_template_file(self, tpl):
        if not os.path.exists(tpl):
            raise FileNotFoundError("not found %s" % tpl)

    def _load_template(self, path):
        self.cheak_template_file(path)  
        self.template = self.template_env.from_string(open(path, "r").read())

    def _load_release_note(self):
        tpl = "./docker/templates/release-note/Note.md.j2"
        if not os.path.exists(tpl):
            raise FileNotFoundError("not found %s" % tpl)
        log.info("load success Note.md j2 file.")
        self.release_note = self.template_env.from_string(open(tpl, "r").read())

    def targeted(self):
        if self.action == "build":
            log.info("building generating")
            self.generate_matrix_tags()
            self.generate_dockerfile()

        elif self.action == "release":
            log.info("releasing generating...")
            self._load_release_note()
            self.generate_matrix_tags()
            self.generate_release_note()
        else:
            log.error("wrong action %s" % self.action)

    def release(self):
        pass

    def main(self):
        self.targeted()
        log.info("Done")


if __name__ == "__main__":
    Manager.run()
